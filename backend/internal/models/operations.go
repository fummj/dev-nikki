package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"dev_nikki/internal/logger"
)

var (
	notExisitEmailError    = errors.New("this email is not exist")
	failedGetProjectsError = errors.New("failed to get projects")

	emailCount int64
)

// 同じemailが存在しないかをチェック。emailが存在していたらerrorを返す。e=email
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

// ユーザー取得。e=email
func getUser(e string) (*gorm.DB, *User, error) {
	user := &User{Email: e}
	result := DBC.DB.Where("email = ?", e).Take(user)
	if result.Error != nil {
		return result, &User{}, result.Error
	}
	return result, user, nil
}

func GetExistUser(e string) (*User, error) {
	err := IsEmailExist(e)
	if err == nil {
		logger.Slog.Error("does not exist this email", "email", e)
		return &User{}, notExisitEmailError
	}

	_, u, err := getUser(e)
	if err != nil {
		logger.Slog.Error("falied get user from db", "error", err)
		return &User{}, err
	}

	return u, nil
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

// プロジェクト取得。
func GetProjects(id uint) (*gorm.DB, []Project, error) {
	project := []Project{}
	result := DBC.DB.Find(project, DBC.DB.Where("usee_id = ?", id))
	if result.Error != nil {
		logger.Slog.Error("failed to get user's projects", "error", result.Error.Error())
		return result, project, failedGetProjectsError
	}
	return result, project, nil
}
