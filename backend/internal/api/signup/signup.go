package signup

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"dev_nikki/internal/api/response"
	"dev_nikki/internal/authN"
	"dev_nikki/internal/logger"
	"dev_nikki/internal/models"
	"dev_nikki/pkg/utils"
)

const (
	charset   string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltCount int    = 16
)

var (
	signupError = errors.New("すでにこのメールアドレスは存在しています。")

	signupFailedResponse = response.SignUpResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			UserID:   0,
			Username: "",
			ErrorMsg: signupError.Error(),
		},
	}
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

func GetPepper() string {
	p := utils.GetEnv(models.EnvPath)["PEPPER"]
	return p
}

func GenerateSalt() string {
	salt := make([]byte, saltCount)
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < saltCount; i++ {
		r := randSeed.Intn(len(charset))
		salt[i] = charset[r]
	}

	return string(salt)
}

func PasswordHashing(p string, salt string) (string, error) {
	b := salt + p + GetPepper()
	h := sha256.New()
	h.Write([]byte(b))
	s := fmt.Sprintf("%x", string(h.Sum(nil)))

	logger.Slog.Debug("completed password hashing", "hashed", s)

	return s, nil
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
		u.salt = GenerateSalt()
	}

	if p, err := PasswordHashing(u.password, u.salt); err != nil {
		return models.DBC.DB, &models.User{}, err
	} else {
		u.password = p
		return models.CreateUser(u.name, u.email, u.password, u.salt)
	}
}

func SignUp(c echo.Context) error {
	_, u, err := createUser(c)

	if err != nil {
		logger.Slog.Error("Failed to create user", "error", err)
		return c.JSON(http.StatusUnprocessableEntity, signupFailedResponse)
	}

	tokenString, err := authN.GenerateJWT(u)
	if err != nil {
		logger.Slog.Error("Failed to create JWT", "error", err)
		return c.JSON(http.StatusUnprocessableEntity, signupFailedResponse)
	}

	logger.Slog.Info("Success to create user, JWT", "user", u)
	resp := response.SignUpResponse{
		Common: response.CommonResponse{
			Status:   "success signup",
			UserID:   u.ID,
			Username: u.Username,
			ErrorMsg: "",
		},
	}
	authN.SetJWTCookie(c, tokenString)
	return c.JSON(http.StatusOK, resp)
}
