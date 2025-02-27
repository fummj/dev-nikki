package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
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
