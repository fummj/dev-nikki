package authN

import (
	"crypto/sha256"
	"fmt"

	"dev_nikki/internal/logger"
	"dev_nikki/internal/models"
)

func PasswordHashing(p string, salt string) (string, error) {
	b := salt + p + models.GetPepper()
	h := sha256.New()
	h.Write([]byte(b))
	s := fmt.Sprintf("%x", string(h.Sum(nil)))

	logger.Slog.Debug("completed password hashing", "hashed", s)

	return s, nil
}
