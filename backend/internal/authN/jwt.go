package authN

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"dev_nikki/internal/logger"
	"dev_nikki/internal/models"
)

var (
	publicKeyPath        string        = "public_jwt.pem"
	privateKeyPath       string        = "private_jwt.pem"
	errInvalidParseToKey error         = errors.New("does not match the parse format")
	errInvalidJWT        error         = errors.New("this jwt is invalid")
	algorithm            string        = "EdDSA"
	period               time.Duration = time.Hour * 4
	iss                  string        = "dev-nikki"
	sub                  string        = "Accusess Token"
	exp                  time.Time     = time.Now().Add(period)
	iat                  time.Time     = time.Now()
	jti                  string        = uuid.NewString()

	KeysKeeper jwtKeysKeeper = NewJWTKeysKeeper()
)

type CustomClaim struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewClaim(id int, name, email string) CustomClaim {
	return CustomClaim{
		UserID:   id,
		Username: name,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(iat),
			ID:        jti,
		},
	}
}

// 署名する前のトークンを返す。
func CreatePreSignedToken(u CustomClaim) *jwt.Token {
	t := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, u)
	return t
}

// 署名されたJWTを生成して返す。
func createJWT(t *jwt.Token, key jwtKeysKeeper) (string, error) {
	tokenString, err := t.SignedString(key.Priv)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 署名されたJWTを生成して返す。
func GenerateJWT(u *models.User) (string, error) {
	claim := NewClaim(int(u.ID), u.Username, u.Email)
	tokenString, err := createJWT(CreatePreSignedToken(claim), KeysKeeper)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(s string, key ed25519.PublicKey) (*jwt.Token, error) {
	var claims CustomClaim
	// ↓
	t, err := jwt.ParseWithClaims(s, &claims, func(token *jwt.Token) (any, error) {
		if token.Header["alg"] == algorithm {
			err := token.Method.Verify(s, token.Signature, key)
			if err != nil {
				return key, nil
			}
		}
		return ed25519.PublicKey{}, errInvalidJWT
	}, jwt.WithValidMethods([]string{algorithm}), jwt.WithIssuedAt(), jwt.WithIssuer(iss))

	if err != nil {
		return t, err
	}

	if !t.Valid {
		return t, errors.New("invalid token in ParseJWT function")
	}

	logger.Slog.Info("ParseJWT", "claims", t.Claims)
	return t, err
}

type keyLoader interface {
	extractPemData(string) (*pem.Block, error)
	parseToPublicKey(pem *pem.Block) error
	parseToPrivateKey(pem *pem.Block) error
	Load()
}

type jwtKeysKeeper struct {
	privPath string
	publPath string
	Priv     ed25519.PrivateKey
	Publ     ed25519.PublicKey
}

func NewJWTKeysKeeper() jwtKeysKeeper {
	j := jwtKeysKeeper{
		privPath: privateKeyPath,
		publPath: publicKeyPath,
		Priv:     ed25519.PrivateKey{},
		Publ:     ed25519.PublicKey{},
	}
	j.Load()

	return j
}

func (keys jwtKeysKeeper) extractPemData(filepath string) (*pem.Block, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil || len(contents) == 0 {
		return nil, fmt.Errorf("extractPemData func: %s", err)
	}

	pemData, _ := pem.Decode(contents)
	if pemData == nil {
		return nil, fmt.Errorf("this file are not \"pem\" format: %s", filepath)
	}
	return pemData, err
}

func (keys *jwtKeysKeeper) parseToPrivateKey(pem *pem.Block) error {
	k, _ := x509.ParsePKCS8PrivateKey(pem.Bytes)

	priv, ok := k.(ed25519.PrivateKey)
	if !ok {
		return errInvalidParseToKey
	}

	keys.Priv = priv
	return nil
}

func (keys *jwtKeysKeeper) parseToPublicKey(pem *pem.Block) error {
	k, _ := x509.ParsePKIXPublicKey(pem.Bytes)

	publ, ok := k.(ed25519.PublicKey)
	if !ok {
		return errInvalidParseToKey
	}

	keys.Publ = publ
	return nil
}

func (keys *jwtKeysKeeper) Load() {
	privPem, err := keys.extractPemData(keys.privPath)
	if err != nil {
		fmt.Println(err)
	}
	keys.parseToPrivateKey(privPem)

	publPem, err := keys.extractPemData(keys.publPath)
	if err != nil {
		fmt.Println(err)
	}
	keys.parseToPublicKey(publPem)
}
