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
func IsEmailExist(tx *gorm.DB, e string) error {
	tx.Table("users").Count(&emailCount)
	if 0 == emailCount {
		return nil
	}

	var user User
	result := tx.Find(&user, "email = ?", e)
	if 0 == result.RowsAffected {
		return nil
	}
	return errors.New(fmt.Sprintf("Failed: this email(%s) is already exist", e))
}

// ユーザー取得。e=email
func getUser(tx *gorm.DB, e string) (*gorm.DB, *User, error) {
	user := &User{Email: e}
	result := tx.Where("email = ?", e).Take(user)
	if result.Error != nil {
		return result, &User{}, result.Error
	}
	return result, user, nil
}

func GetExistUser(tx *gorm.DB, e string) (*User, error) {
	err := IsEmailExist(tx, e)
	if err == nil {
		logger.Slog.Error("does not exist this email", "email", e)
		return &User{}, notExisitEmailError
	}

	_, u, err := getUser(tx, e)
	if err != nil {
		logger.Slog.Error("falied get user from db", "error", err)
		return &User{}, err
	}

	return u, nil
}

// ユーザー作成。n=name, e=email, p=password, s=salt
func CreateUser(tx *gorm.DB, n, e, p, s string) (*gorm.DB, *User, error) {
	user := &User{
		Username: n,
		Email:    e,
		Password: p,
		Salt:     s,
	}
	result := tx.Create(user)
	if result.Error != nil {
		return result, &User{}, result.Error
	}
	return result, user, nil
}

// ユーザーに紐づくプロジェクトに引数nと同じ名前のプロジェクトが存在するか調べる。
func isExistProject(tx *gorm.DB, n string, uid uint) error {
	var project Project
	result := tx.Where(&Project{Name: n, UserID: uid}, "name", "user_id").Find(&project)

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

// プロジェクト作成。n=name, d=Description, uid=user_id
func CreateProject(tx *gorm.DB, n, d string, uid uint) (*gorm.DB, *Project, error) {
	if err := isExistProject(tx, n, uid); err != nil {
		return &gorm.DB{}, &Project{}, err
	}

	project := &Project{
		Name:        n,
		Description: d,
		UserID:      uid,
	}
	result := tx.Create(project)
	if result.Error != nil {
		return result, &Project{}, result.Error
	}
	return result, project, nil
}

// project_idに紐づくプロジェクト取得
func GetProject(tx *gorm.DB, id uint) (*gorm.DB, Project, error) {
	var project Project
	result := tx.First(&project, id)
	if result.Error != nil {
		logger.Slog.Error("failed to get project from user", "error", result.Error.Error())
		return result, project, failedGetProjectError
	}
	return result, project, nil
}

// user_idに紐づく複数のプロジェクト取得。
func GetProjects(tx *gorm.DB, uid uint) (*gorm.DB, []Project, error) {
	var projects []Project
	result := tx.Where(&Project{UserID: uid}, "user_id").Find(&projects)
	if result.Error != nil {
		logger.Slog.Error("failed to get projects from user", "error", result.Error.Error())
		return result, projects, failedGetProjectsError
	}
	return result, projects, nil
}

// フォルダ取得。
func GetFolders(tx *gorm.DB, userID uint, projectID uint) (*gorm.DB, []Folder, error) {
	var folders []Folder
	result := tx.Where(&Folder{UserID: userID, ProjectID: projectID},
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
func GetFiles(tx *gorm.DB, userID, projectID, folderID uint) ([]File, error) {
	var files []File
	result := tx.Where(&File{UserID: userID, ProjectID: projectID, FolderID: &folderID},
		"user_id",
		"project_id",
		"folder_id",
	).Find(&files)

	if result.Error != nil {
		logger.Slog.Error("failed to get files from user's project's folder", "error", result.Error.Error())
		return files, failedGetFilesError
	}

	return files, nil
}

// フォルダに関連していないファイル取得。
func GetNoFolderFiles(tx *gorm.DB, userID, projectID uint) ([]File, error) {
	var files []File

	result := tx.Where(
		&File{UserID: userID, ProjectID: projectID, FolderID: nil},
		"user_id",
		"project_id",
		"folder_id",
	).Find(&files)

	if result.Error != nil {
		logger.Slog.Error(result.Error.Error())
		return files, result.Error
	}

	return files, nil
}

// 一つのファイルを取得。
func GetFile(tx *gorm.DB, fileID uint) (File, error) {
	var file File
	result := tx.First(&file, fileID)

	if result.Error != nil {
		logger.Slog.Error("", "error", result.Error.Error())
		return file, result.Error
	}

	return file, nil
}

// ファイルのContentを更新。
func UpdateFile(tx *gorm.DB, fileID uint, content string) error {
	result := tx.Model(&File{ID: fileID}).Update("content", content)
	if result.Error != nil {
		logger.Slog.Error("failed to update file content", "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// フォルダ作成
func CreateFolder(tx *gorm.DB, n string, ui, pi uint, pfi *uint) (Folder, error) {
	fo := Folder{
		Name:           n,
		UserID:         ui,
		ProjectID:      pi,
		ParentFolderID: pfi,
	}

	result := tx.Create(&fo)
	if result.Error != nil {
		logger.Slog.Error(result.Error.Error())
		return fo, result.Error
	}

	return fo, nil
}

// ファイル作成
func CreateFile(tx *gorm.DB, n string, ui, pi uint, fi *uint) (File, error) {
	f := File{
		Name:      n,
		UserID:    ui,
		ProjectID: pi,
		FolderID:  fi,
	}

	result := tx.Create(&f)
	if result.Error != nil {
		logger.Slog.Error(result.Error.Error())
		return f, result.Error
	}

	return f, nil
}

// フォルダ削除
func DeleteFolder(tx *gorm.DB, foi uint) error {
	fo := Folder{
		ID: foi,
	}

	result := tx.Delete(&fo)
	if result.Error != nil {
		logger.Slog.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// ファイル削除
func DeleteFile(tx *gorm.DB, fi uint) error {
	f := File{
		ID: fi,
	}

	result := tx.Delete(&f)
	if result.Error != nil {
		logger.Slog.Error(result.Error.Error())
		return result.Error
	}

	return nil
}
