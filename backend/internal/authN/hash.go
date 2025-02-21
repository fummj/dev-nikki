package authN

import (
	"crypto/sha256"
	"fmt"
	"log/slog"

	"dev_nikki/internal/models"
)

func PasswordHashing(p string, salt string) (string, error) {
	b := salt + p + models.GetPepper()
	h := sha256.New()
	h.Write([]byte(b))
	s := fmt.Sprintf("%x", string(h.Sum(nil)))

	slog.Debug("completed password hashing", "hashed", s)

	return s, nil
}
