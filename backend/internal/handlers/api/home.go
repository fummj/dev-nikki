package api

import (
	"github.com/labstack/echo/v4"

	"dev_nikki/internal/api/home"
)

type HomeHandler struct {
}

func (h HomeHandler) Home(c echo.Context) error {
	return home.Home(c)
}

func (h HomeHandler) PreHome(c echo.Context) error {
	return home.PreHome(c)
}

func (h HomeHandler) PostPreHome(c echo.Context) error {
	return home.PostPreHome(c)
}
