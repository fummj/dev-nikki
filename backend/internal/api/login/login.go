package login

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type loginFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ReturnFormData(c echo.Context) error {
	email := c.Request().FormValue("email")
	password := c.Request().FormValue("password")
	d := &loginFormData{email, password}
	fmt.Println("login-form data: ", d)
	return c.JSON(http.StatusOK, d)
}
