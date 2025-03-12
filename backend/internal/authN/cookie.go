package authN

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"dev_nikki/internal/logger"
)

// cookieにJWTを保存。
func SetJWTCookie(c any, token string) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	echoContext, ok := c.(echo.Context)
	if ok {
		fmt.Println("SetJWTCookie cookie(echo.Context)", cookie)
		echoContext.SetCookie(cookie)
		return
	}

	req, ok := c.(*http.Request)
	if ok {
		fmt.Println("SetJWTCookie cookie(http.Request)", cookie)
		req.AddCookie(cookie)
		return
	}

}

// cookieのJWTを検証
func ParseJWTCookie(c echo.Context) (*jwt.Token, error) {
	cookie, err := c.Cookie("access_token")
	if err != nil {
		logger.Slog.Error("Failed to get jwt from cookie", "error", err)
		return jwt.New(&jwt.SigningMethodEd25519{}), err
	}

	t, err := ParseJWT(cookie.Value, KeysKeeper.Publ)
	fmt.Println("cookie.Value", cookie.Value)
	if err != nil {
		logger.Slog.Error("Failed to parse jwt from cookie", "err", err)
		return t, err
	}

	logger.Slog.Info("Success to parse jwt from cookie", "jwt-token", t)
	return t, nil
}
