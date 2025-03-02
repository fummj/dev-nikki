package api

import (
	"github.com/labstack/echo/v4"

	"dev_nikki/internal/api/login"
	"dev_nikki/internal/api/signup"
)

type AuthHandler struct {
}

func (a AuthHandler) SignUp(c echo.Context) error {
	return signup.SignUp(c)
}

func (a AuthHandler) Login(c echo.Context) error {
	return login.Login(c)
}
