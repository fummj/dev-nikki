package signup

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"dev_nikki/internal/authN"
	"dev_nikki/internal/models"
)

var (
	userMap = map[string]string{
		"name":     "",
		"email":    "",
		"password": "",
		"salt":     "",
	}
)

type userData struct {
	name     string
	email    string
	password string
	salt     string
}

type responseData struct {
	Status   string `json:"status"`
	ID       uint   `json:"id"`
	JWT      string `json:"jwt"`
	ErrorMsg string `json:"errorMsg"`
}

func newUserData(m map[string]string) *userData {
	return &userData{
		name:     m["name"],
		email:    m["email"],
		password: m["password"],
		salt:     m["salt"],
	}
}

// バリデーションをした後に送られてきたデータを元にUserを作成してDBに保存する。
func createUser(c echo.Context) (*gorm.DB, *models.User, error) {
	for k := range userMap {
		userMap[k] = c.Request().FormValue(k)
	}
	u := newUserData(userMap)
	slog.Info("check request form data", "name", u.name, "email", u.email, "pass", u.password)

	if err := models.IsEmailExist(u.email); err != nil {
		slog.Error("this email is already exist")
		return models.DBC.DB, &models.User{}, err
	}

	if err := authN.Validation(u.email, u.password); err != nil {
		return models.DBC.DB, &models.User{}, err
	} else {
		u.salt = models.GenerateSalt()
	}

	if p, err := authN.PasswordHashing(u.password, u.salt); err != nil {
		return models.DBC.DB, &models.User{}, err
	} else {
		u.password = p
		return models.CreateUser(u.name, u.email, u.password, u.salt)
	}
}

func SendUserData(c echo.Context) error {
	_, user, err := createUser(c)

	if err != nil {
		slog.Error("Failed to create user", "error", err)
		resp := responseData{Status: "failed", ErrorMsg: err.Error()}
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	slog.Info("Success to create user", "user", user)
	resp := responseData{"success", user.ID, "jwt", ""}
	return c.JSON(http.StatusOK, resp)
}
