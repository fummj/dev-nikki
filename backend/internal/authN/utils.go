package authN

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"dev_nikki/internal/models"
)

const (
	charset        string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltCount      int    = 16
	minPasswordLen int    = 8
	maxPasswordLen int    = 100
)

var count int64

type userProcessError struct {
	file string
	line string
	msg  string
}

// パスワードのバリデーションチェック
func PasswordValidate(p string) bool {
	// 8文字以下or10文字以上かチェック(英数字)
	s := []rune(p)
	if len(s) < minPasswordLen || (len(s) > maxPasswordLen) {
		slog.Error("this password does not meet the requirements")
		return false
	}
	// 大文字、小文字、数字がそれぞれ1つあるかをチェック
	u, l, n := false, false, false

	for i := 0; i < len(p); i++ {
		// 数字の存在チェック
		if _, err := strconv.Atoi(string(s[i])); err == nil {
			n = true
			continue
		}
		// 大文字の存在チェック
		if strings.ToUpper(string(s[i])) == string(s[i]) {
			slog.Debug(strings.ToUpper(string(s[i])), string(s[i]), 1)
			u = true
			continue
		}
		// 小文字の存在チェック
		if strings.ToLower(string(s[i])) == string(s[i]) {
			slog.Debug(strings.ToLower(string(s[i])), string(s[i]), 2)
			l = true
			continue
		}
	}

	if u && l && n {
		return true
	}
	return false
}

func PasswordHashing(p string, salt string) (string, error) {
	if !PasswordValidate(p) {
		return "", errors.New("password validation failed")
	}
	b := salt + p + models.GetPepper()
	h := sha256.New()
	h.Write([]byte(b))
	s := fmt.Sprintf("%x", string(h.Sum(nil)))

	slog.Debug("completed password hashing", "hashed", s)

	return s, nil
}
