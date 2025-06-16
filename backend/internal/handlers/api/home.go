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

func (h HomeHandler) UpdateMarkdown(c echo.Context) error {
	return home.UpdateMarkdown(c)
}

func (h HomeHandler) CreateNewFolder(c echo.Context) error {
	return home.CreateNewFolder(c)
}

func (h HomeHandler) CreateNewFile(c echo.Context) error {
	return home.CreateNewFile(c)
}

func (h HomeHandler) DeleteFolder(c echo.Context) error {
	return home.DeleteFolder(c)
}

func (h HomeHandler) DeleteFile(c echo.Context) error {
	return home.DeleteFile(c)
}
