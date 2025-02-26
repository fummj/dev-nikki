package signup

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"dev_nikki/internal/authN"
	"dev_nikki/internal/logger"
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

type signupResponseData struct {
	Status   string `json:"status"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
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
	logger.Slog.Info("check request form data", "name", u.name, "email", u.email, "pass", u.password)

	if err := models.IsEmailExist(u.email); err != nil {
		logger.Slog.Error("this email is already exist")
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

// 署名されたJWTを生成する。
func generateJWT(u *models.User) (string, error) {
	claim := authN.NewClaim(int(u.ID), u.Username, u.Email)
	tokenString, err := authN.CreateJWT(authN.CreatePreSignedToken(claim), authN.KeysKeeper)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SignUp(c echo.Context) error {
	_, user, err := createUser(c)

	if err != nil {
		logger.Slog.Error("Failed to create user", "error", err)
		resp := signupResponseData{Status: "failed", ErrorMsg: err.Error()}
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	tokenString, err := generateJWT(user)
	if err != nil {
		logger.Slog.Error("Failed to create JWT", "error", err)
		resp := signupResponseData{Status: "failed", ErrorMsg: err.Error()}
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	logger.Slog.Info("Success to create user, JWT", "user", user)
	resp := signupResponseData{"success", user.ID, user.Username, ""}
	authN.SetJWTCookie(c, tokenString)
	return c.JSON(http.StatusOK, resp)
}
