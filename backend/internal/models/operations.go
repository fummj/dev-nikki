package models

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"dev_nikki/pkg/utils"
)

const (
	charset   string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltCount int    = 16
)

var emailCount int64

// 同じemailが存在しないかをチェック
func IsEmailExist(e string) error {
	DBC.DB.Table("users").Count(&emailCount)
	if 0 == emailCount {
		return nil
	}

	var user User
	result := DBC.DB.Find(&user, "email = ?", e)
	if 0 == result.RowsAffected {
		return nil
	}
	return errors.New(fmt.Sprintf("Failed: this email(%s) is already exist", e))
}

func GetPepper() string {
	p := utils.GetEnv(EnvPath)["PEPPER"]
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

// ユーザー作成。n=name, e=email, p=password, s=salt
func CreateUser(n, e, p, s string) (*gorm.DB, *User, error) {
	user := &User{
		Username: n,
		Email:    e,
		Password: p,
		Salt:     s,
	}
	result := DBC.DB.Create(user)
	if result.Error != nil {
		return result, &User{}, result.Error
	}
	return result, user, nil
}
