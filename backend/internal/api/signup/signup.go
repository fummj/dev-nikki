package signup

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/authN"
	"dev_nikki/internal/models"
)

type signupFormData struct {
	Status   string `json:"status"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ReturnCreatedUserData(c echo.Context) error {
	name := c.Request().FormValue("name")
	email := c.Request().FormValue("email")
	password := c.Request().FormValue("password")

	userMap := map[string]string{
		"username": name,
		"email":    email,
		"password": password,
		"salt":     authN.GenerateSalt(),
	}

	fmt.Println("start!")
	result, user, err := models.CreateUser(models.DBC.DB, userMap)

	if err != nil {
		d := &signupFormData{name, email, password, "failed"}
		fmt.Println("signup-form data: ", d)
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, d)
	}

	d := signupFormData{name, email, password, "success!"}
	fmt.Println(result, user, d)
	return c.JSON(http.StatusOK, user)
}
