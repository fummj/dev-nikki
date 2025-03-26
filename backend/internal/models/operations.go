package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"dev_nikki/internal/logger"
)

var (
	notExisitEmailError      = errors.New("this email is not exist")
	failedGetProjectError    = errors.New("failed to get project")
	failedGetProjectsError   = errors.New("failed to get projects")
	failedGetFoldersError    = errors.New("failed to get folders")
	failedGetFilesError      = errors.New("failed to get files")
	AlreadyExistProjectError = errors.New("プロジェクト名が重複しています。")

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

// ユーザーに紐づくプロジェクトに引数nと同じ名前のプロジェクトが存在するか調べる。
func isExistProject(n string, id uint) error {
	var project Project
	result := DBC.DB.Where(&Project{Name: n, UserID: id}, "name", "user_id").Find(&project)

	if result.Error != nil {
		logger.Slog.Error(result.Error.Error())
		return result.Error
	}

	if result.RowsAffected != 0 {
		logger.Slog.Error(AlreadyExistProjectError.Error())
		return AlreadyExistProjectError
	}

	return nil
}

// プロジェクト作成。n=name, d=Description, id=user_id
func CreateProject(n, d string, id uint) (*gorm.DB, *Project, error) {
	if err := isExistProject(n, id); err != nil {
		return &gorm.DB{}, &Project{}, err
	}

	project := &Project{
		Name:        n,
		Description: d,
		UserID:      id,
	}
	result := DBC.DB.Create(project)
	if result.Error != nil {
		return result, &Project{}, result.Error
	}
	return result, project, nil
}

// user_idに紐づくプロジェクト取得
func GetProject(id uint) (*gorm.DB, Project, error) {
	var project Project
	result := DBC.DB.Where(&Project{UserID: id}, "user_id").Find(&project)
	if result.Error != nil {
		logger.Slog.Error("failed to get project from user", "error", result.Error.Error())
		return result, project, failedGetProjectError
	}
	return result, project, nil
}

// user_idに紐づく複数のプロジェクト取得。
func GetProjects(id uint) (*gorm.DB, []Project, error) {
	var projects []Project
	result := DBC.DB.Where(&Project{UserID: id}, "user_id").Find(&projects)
	if result.Error != nil {
		logger.Slog.Error("failed to get projects from user", "error", result.Error.Error())
		return result, projects, failedGetProjectsError
	}
	return result, projects, nil
}

// フォルダ取得。
func GetFolders(userID uint, projectID uint) (*gorm.DB, []Folder, error) {
	var folders []Folder
	result := DBC.DB.Where(&Folder{UserID: userID, ProjectID: projectID},
		"user_id",
		"project_id",
	).Find(&folders)

	if result.Error != nil {
		logger.Slog.Error("failed to get folder from user's project", "error", result.Error.Error())
		return result, folders, failedGetFoldersError
	}
	return result, folders, nil
}

// ファイル取得。
func GetFiles(userID, projectID, folderID uint) (*gorm.DB, []File, error) {
	var files []File
	result := DBC.DB.Where(&File{UserID: userID, ProjectID: projectID, FolderID: folderID},
		"user_id",
		"project_id",
		"folder_id",
	).Find(&files)

	if result.Error != nil {
		logger.Slog.Error("failed to get files from user's project's folder", "error", result.Error.Error())
		return result, files, failedGetFilesError
	}
	return result, files, nil
}
