package login

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/api/response"
	"dev_nikki/internal/api/signup"
	"dev_nikki/internal/authN"
	"dev_nikki/internal/logger"
	"dev_nikki/internal/models"
)

var (
	noMatchPasswordError = errors.New("password do not match")
	loginError           = errors.New("メールアドレスかパスワードが正しくありません。")

	loginFailedResponse = response.LoginResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			UserID:   0,
			Username: "",
			Email:    "",
			ErrorMsg: loginError.Error(),
		},
	}
)

// requestのpassword(hashed)がDB内のものと一致するかの検証。 p=requestのpassword
func verifyHashedPassword(p string, u *models.User) error {
	// requestのpasswordをhashed
	hashed, err := signup.PasswordHashing(p, u.Salt)
	if err != nil {
		logger.Slog.Error("login failed", "error", err)
		return err
	}

	// requestのpassword(hashed)がDB内のものと一致するかの検証
	if hashed != u.Password {
		return noMatchPasswordError
	}

	return nil
}

func Login(c echo.Context) error {
	email := c.Request().FormValue("email")
	password := c.Request().FormValue("password")

	u, err := models.GetExistUser(models.DBC.DB, email)
	if err != nil {
		// front側に具体的なエラー内容は流さないようにloginFailedResponseを使用。
		return c.JSON(http.StatusUnauthorized, loginFailedResponse)
	}

	err = verifyHashedPassword(password, u)
	if err != nil {
		// front側に具体的なエラー内容は流さないようにloginErrorを使用。
		return c.JSON(http.StatusUnauthorized, loginFailedResponse)
	}

	tokenString, err := authN.GenerateJWT(u)
	if err != nil {
		logger.Slog.Error("Failed to create JWT", "error", err)
		return c.JSON(http.StatusUnauthorized, loginFailedResponse)
	}

	authN.SetJWTCookie(c, tokenString)
	resp := response.LoginResponse{
		Common: response.CommonResponse{
			Status:   "success login",
			UserID:   u.ID,
			Username: u.Username,
			Email:    u.Email,
			ErrorMsg: "",
		},
	}

	return c.JSON(http.StatusOK, resp)
}
