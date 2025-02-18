package models

import (
	"errors"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"dev_nikki/pkg/utils"
)

const (
	charset   string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltCount int    = 16
)

var count int64

// 同じemailが存在しないかをチェック
func IsEmailExist(db *gorm.DB, e string) error {
	db.Table("users").Count(&count)
	if 0 == count {
		return nil
	}

	var user User
	result := db.Find(&user, "email = ?", e)
	if 0 == result.RowsAffected {
		return errors.New("")
	}
	return nil
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

// ユーザー作成
func CreateUser(db *gorm.DB, userMap map[string]string) (*gorm.DB, User, error) {
	// emailが一意なので存在していないかのチェック
	// 残すは同じemailが既存在しないかをチェックするだけ。1/5
	// p, err := PasswordHashing(userMap["password"], userMap["salt"])
	// if err != nil {
	// 	return db, User{}, err
	// }

	// もしも他に同じemailが存在していたらここでエラー吐き出して終わらせる。
	// err = IsEmailExist(db, userMap["email"])
	// if err != nil {
	// 	return db, User{}, err
	// }

	user := &User{
		Username: userMap["username"],
		Email:    userMap["email"],
		// Password: p,
		Password: userMap["password"],
		Salt:     userMap["salt"],
	}
	result := db.Create(user)
	if result.Error != nil {
		err := errors.New(result.Error.Error())
		return result, User{}, err
	}
	return result, *user, nil
}
