package authN

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"

	"dev_nikki/internal/logger"
	"dev_nikki/internal/models"
	"dev_nikki/pkg/utils"
)

const (
	googleURL          = "https://accounts.google.com"
	redirectURL        = "http://localhost:8080/auth/callback"
	successRedirectURL = "http://localhost:8080/prehome"
	successRequestURL  = "http://localhost:8080/api/home/prehome"
)

var (
	notExistProviderError           = errors.New("this provider is not exist")
	failedGetAuthorizationCodeError = errors.New("failed to get authorization code")
	notMatchStateError              = errors.New("response state does not match request state")
	failedExtractIDTokenError       = errors.New("failed to extract id_token")
	failedVerifierIDTokenError      = errors.New("failed to verifier id_token")

	envPath     = ".env"
	credentials = map[string]string{
		"CLIENT_ID":     "",
		"CLIENT_SECRET": "",
	}
	state    string
	verifier string
	provider *oidc.Provider
)

func init() {
	getCredentials()
}

func getCredentials() {
	m := utils.GetEnv(envPath)
	for k := range credentials {
		credentials[k] = m[k]
	}
}

func generateState() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	state = base64.RawURLEncoding.EncodeToString(b)
	return state
}

func afterSuccessedOAuth2(c echo.Context, t string) ([]byte, error) {
	// successRequestURLのリクエスト用
	req, err := http.NewRequest("GET", successRequestURL, http.NoBody)
	if err != nil {
		return []byte{}, err
	}
	SetJWTCookie(req, t)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	c.Redirect(http.StatusFound, successRedirectURL)

	logger.Slog.Info("success OAuth2.0(OIDC)")
	return resBody, nil
}

func oauth2SignUp(c echo.Context, n, e string) error {
	_, u, err := models.CreateUser(n, e, "", "")
	if err != nil {
		return err
	}
	tokenString, err := GenerateJWT(u)
	// レスポンス用
	SetJWTCookie(c, tokenString)

	resBody, err := afterSuccessedOAuth2(c, tokenString)
	if err != nil {
		return err
	}

	logger.Slog.Info("success: SignUp")
	return c.JSON(http.StatusOK, resBody)
}

func oauth2Login(c echo.Context, e string) error {
	u, err := models.GetExistUser(e)
	if err != nil {
		return err
	}

	tokenString, err := GenerateJWT(u)
	SetJWTCookie(c, tokenString)

	resBody, err := afterSuccessedOAuth2(c, tokenString)
	if err != nil {
		return err
	}

	logger.Slog.Info("success: Login")
	return c.JSON(http.StatusOK, resBody)
}

func newOAuth2Config(c echo.Context) (*oauth2.Config, error) {
	provider, _ = oidc.NewProvider(c.Request().Context(), googleURL)

	oauth2Config := &oauth2.Config{
		ClientID:     credentials["CLIENT_ID"],
		ClientSecret: credentials["CLIENT_SECRET"],
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return oauth2Config, nil
}

func OAuth2(c echo.Context) error {
	oauth2Config, _ := newOAuth2Config(c)

	verifier = oauth2.GenerateVerifier()
	url := oauth2Config.AuthCodeURL(generateState(), oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("the URL that authorization server returns to the user: %v\n", url)

	http.Redirect(c.Response().Writer, c.Request(), url, http.StatusFound)

	return nil
}

func OAuth2Callback(c echo.Context) error {
	// stateが同じかをチェック。
	if state != c.FormValue("state") {
		return notMatchStateError
	}

	code := c.FormValue("code")
	if code == "" {
		return failedGetAuthorizationCodeError
	}

	oauth2Config, _ := newOAuth2Config(c)
	oauthToken, err := oauth2Config.Exchange(c.Request().Context(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		return err
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return failedExtractIDTokenError
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: credentials["CLIENT_ID"]})
	idToken, err := idTokenVerifier.Verify(c.Request().Context(), rawIDToken)
	if err != nil {
		return err
	}

	var claims struct {
		Name     string
		Email    string
		Verifier bool
	}

	if err = idToken.Claims(&claims); err != nil {
		return err
	}

	logger.Slog.Debug("result claims", "claims", claims)

	email := claims.Email

	// emailが存在してたらLogin, 存在してなかったらSignUp。
	err = models.IsEmailExist(email)
	if err != nil {
		logger.Slog.Info("OAuth2(OIDC) Login")
		oauth2Login(c, email)
		return nil
	}

	logger.Slog.Info("OAuth2(OIDC) SignUp")
	oauth2SignUp(c, claims.Name, claims.Email)

	return nil
}
