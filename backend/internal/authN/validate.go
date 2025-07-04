package authN

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"

	"dev_nikki/internal/logger"
)

var (
	emailValidationError         error = errors.New("this email is invalid")
	passwordValidationError      error = errors.New("this password is invalid")
	verifyJWTAgainstRequestError error = errors.New("some of the information in claims differs from some of the information requested")
)

const (
	emailPattern   string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	maxEmailLen    int    = 254
	minPasswordLen int    = 8
	maxPasswordLen int    = 100
)

var passwordPatterns []string = []string{
	`\d`,
	`[a-z]`,
	`[A-Z]`,
	"[!@#$%^&*()_+[\\]{};':\"\\|,.<>?/`~-]",
}

// emailのバリデーションチェック
func EmailValidation(e string) error {
	if utf8.RuneCountInString(e) > maxEmailLen {
		logger.Slog.Error("this email is too long", "email", e)
		return emailValidationError
	}

	emailRegex := regexp.MustCompile(emailPattern)
	if result := emailRegex.MatchString(e); !result {
		return emailValidationError
	}
	return nil
}

// passwordのバリデーションチェック
func PasswordValidation(p string) error {
	count := utf8.RuneCountInString(p)
	if count < minPasswordLen || count > maxPasswordLen {
		logger.Slog.Error("this password is wrong", "password", p)
		return passwordValidationError
	}

	// 大文字、小文字、数字、記号が1つ以上あるかのチェック
	for _, pv := range passwordPatterns {
		passwordRegex := regexp.MustCompile(pv)
		if result := passwordRegex.MatchString(p); !result {
			return passwordValidationError
		}
	}
	return nil
}

// email、passwordのバリデーションチェックを行う。e=email、p=password
func Validation(e string, p string) error {
	if err := EmailValidation(e); err != nil {
		fmt.Println(err)
		return err
	}

	if err := PasswordValidation(p); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// cid(claimsのuser_id), rbid(リクエストボディから取得したuser_id)が一致するかのチェック
func VerifyJWTAgainstRequest(cid, rbid uint) error {

	if cid != rbid {
		logger.Slog.Error("access is denied because claims information does not match the request-user information", "claims-userid", cid, "request-userid", rbid)
		return verifyJWTAgainstRequestError
	}

	return nil
}
