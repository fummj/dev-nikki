package authN

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"dev_nikki/internal/models"
)

var (
	userID        uint   = 1
	test_username string = "test_jwt"
	test_email    string = "test@jwt.com"
	test_password string = "398asdlkfj$Kad"
	test_salt     string = "aeijae3$d.I9"
)

func setUpKeysKeeper() jwtKeysKeeper {
	return NewJWTKeysKeeper()
}

func setUpCreateClaims(t *testing.T) CustomClaims {
	claims := NewClaim(userID, test_username, test_email)
	if _, ok := interface{}(claims).(CustomClaims); !ok {
		t.Errorf("生成されたものは %Tではなかった。", CustomClaims{})
	}
	return claims
}

func setUpUserData(t *testing.T, tx *gorm.DB) *models.User {
	_, user, err := models.CreateUser(tx, test_username, test_email, test_password, test_salt)
	if err != nil {
		t.Error(err)
	}
	return user
}

func TestNewClaim(t *testing.T) {
	setUpCreateClaims(t)
}

func TestCreatePreSignedToken(t *testing.T) {
	claims := setUpCreateClaims(t)
	preSignedToken := CreatePreSignedToken(claims)

	if sm := (jwt.SigningMethodEd25519{}); sm.Alg() != preSignedToken.Method.Alg() {
		t.Error("署名アルゴリズムが異なる。")
	}

	pst, ok := preSignedToken.Claims.(CustomClaims)
	if !ok {
		t.Errorf("使用されているclaimsが%Tではない。", CustomClaims{})
	}

	if claims.UserID != pst.UserID {
		t.Error("使用したclaimsのデータと生成された署名前段階のトークンのデータが異なる。")
	}
}

func Test_createJWT(t *testing.T) {
	claims := setUpCreateClaims(t)
	preSignedToken := CreatePreSignedToken(claims)
	keysKeeper := setUpKeysKeeper()

	_, err := createJWT(preSignedToken, keysKeeper)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateJWT(t *testing.T) {
	tx := models.DBC.DB.Begin()
	u := setUpUserData(t, tx)

	s, err := GenerateJWT(u)
	if err != nil {
		t.Error(err)
	}

	if s == "" {
		t.Error("署名されたトークンが生成されていない。")
	}
	tx.Rollback()
}

func TestParseJWT(t *testing.T) {
	tx := models.DBC.DB.Begin()
	u := setUpUserData(t, tx)
	k := setUpKeysKeeper()

	s, err := GenerateJWT(u)
	if err != nil {
		t.Error(err)
	}

	tk, err := ParseJWT(s, k.Publ)
	if err != nil {
		t.Error(err)
	}

	if !tk.Valid {
		t.Error("正しいトークンではありません。")
	}
	tx.Rollback()
}

func Test_extractCustomClaims(t *testing.T) {
	tx := models.DBC.DB.Begin()
	u := setUpUserData(t, tx)
	k := setUpKeysKeeper()

	s, err := GenerateJWT(u)
	if err != nil {
		t.Error(err)
	}

	tk, err := ParseJWT(s, k.Publ)
	if err != nil {
		t.Error(err)
	}

	if !tk.Valid {
		t.Error("正しいトークンではありません。")
	}

	_, err = extractCustomClaims(tk)
	if err != nil {
		t.Error(err)
	}
	tx.Rollback()
}

func TestGetExtractedCustomClaims(t *testing.T) {
	tx := models.DBC.DB.Begin()
	u := setUpUserData(t, tx)
	k := setUpKeysKeeper()

	ts, err := GenerateJWT(u)
	if err != nil {
		t.Error(err)
	}

	tk, err := ParseJWT(ts, k.Publ)
	if err != nil {
		t.Error(err)
	}

	if !tk.Valid {
		t.Error("正しいトークンではありません。")
	}

	claims, _ := tk.Claims.(*CustomClaims)

	app := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/home", nil)
	req.Header.Add("Cookie", fmt.Sprintf("access_token=%s", ts))
	rec := httptest.NewRecorder()
	ec := app.NewContext(req, rec)

	extractedClaims, err := GetExtractedCustomClaims(ec)
	if err != nil {
		t.Error(err)
	}

	if extractedClaims.UserID != claims.UserID {
		t.Error("cookieから取得したclaimsの情報と元のデータが異なる。")
	}

	tx.Rollback()
}

func TestNewJWTKeysKeeper(t *testing.T) {
	k := setUpKeysKeeper()
	if _, ok := interface{}(k).(jwtKeysKeeper); !ok {
		t.Error("生成されたものがjwtKeysKeeperではない。")
	}
}

func Test_extractPemData(t *testing.T) {
	k := setUpKeysKeeper()
	if _, err := k.extractPemData(k.privPath); err != nil {
		t.Error(err)
	}

	if _, err := k.extractPemData(k.publPath); err != nil {
		t.Error(err)
	}
}

func Test_parseToPrivateKey(t *testing.T) {
	k := setUpKeysKeeper()
	privPem, err := k.extractPemData(k.privPath)
	if err != nil {
		t.Error(err)
	}

	if err = k.parseToPrivateKey(privPem); err != nil {
		t.Error(err)
	}
}

func Test_parseToPublicKey(t *testing.T) {
	k := setUpKeysKeeper()
	publPem, err := k.extractPemData(k.publPath)
	if err != nil {
		t.Error(err)
	}

	if err = k.parseToPublicKey(publPem); err != nil {
		t.Error(err)
	}
}

func TestLoad(t *testing.T) {
	k := setUpKeysKeeper()
	privPem, err := k.extractPemData(k.privPath)
	if err != nil {
		t.Error(err)
	}

	publPem, err := k.extractPemData(k.publPath)
	if err != nil {
		t.Error(err)
	}

	if err = k.parseToPrivateKey(privPem); err != nil {
		t.Error(err)
	}

	if err = k.parseToPublicKey(publPem); err != nil {
		t.Error(err)
	}
}
