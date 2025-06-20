package models

import (
	"testing"

	"gorm.io/gorm"
)

var (
	username string = "test_user2"
	email    string = "test@test2.com"
	password string = "adaefijfe39adkj"
	// TODO: saltに関してはソースコード見て判断
	salt string = "afijefi9309afj:"
	user *User  = &User{Username: username, Email: email, Password: password, Salt: salt}
)

func setUpUser(t *testing.T, tx *gorm.DB) *User {

	_, user, err := CreateUser(tx, username, email, password, salt)
	if err != nil {
		t.Error(err)
	}
	return user
}

func setUpUntilProject(t *testing.T, tx *gorm.DB) (*User, *Project) {
	setUpUser(t, tx)

	_, project, err := CreateProject(tx, TestProjectName, TestProjectDescription, user.ID)
	if err != nil {
		t.Error(err)
	}
	return user, project
}

func setUpUntilFolder(t *testing.T, tx *gorm.DB) (*User, *Project, Folder) {
	setUpUser(t, tx)

	_, project, err := CreateProject(tx, TestProjectName, TestProjectDescription, user.ID)
	if err != nil {
		t.Error(err)
	}

	folder, err := CreateFolder(tx, TestFolderName, user.ID, project.ID, nil)
	if err != nil {
		t.Error(err)
	}

	if folder.ParentFolderID != nil {
		t.Error("作成されたフォルダが意図していないフォルダに紐づけられている。")
	}

	return user, project, folder
}

func setUpUntilFile(t *testing.T, tx *gorm.DB) (*User, *Project, *Folder, *File) {
	user, project, folder := setUpUntilFolder(t, tx)
	file, err := CreateFile(tx, TestFileName, user.ID, project.ID, &folder.ID)
	if err != nil {
		t.Error(err)
	}

	if *file.FolderID != folder.ID {
		t.Error("作成されたファイルがフォルダに紐付けられていない。")
	}

	return user, project, &folder, &file
}

func TestIsEmailExist(t *testing.T) {
	tx := TestDBC.DB.Begin()

	setUpUser(t, tx)

	err := IsEmailExist(tx, email)
	if err == nil {
		t.Errorf("%sが確認できない。", email)
	}

	err = IsEmailExist(tx, "notexist@email.com")
	if err != nil {
		t.Errorf("%sが存在していないはずなのに検出されている。。", TestUserEmail)
	}

	tx.Rollback()
}

func Test_getUser(t *testing.T) {
	tx := TestDBC.DB.Begin()

	setUpUser(t, tx)

	_, _, err := getUser(tx, email)
	if err != nil {
		t.Error(err)
	}

	tx.Rollback()
}

func TestGetExistUser(t *testing.T) {
	tx := TestDBC.DB.Begin()

	setUpUser(t, tx)

	_, err := GetExistUser(tx, email)
	if err != nil {
		t.Error(err)
	}

	tx.Rollback()
}

func TestCreateUser(t *testing.T) {
	tx := TestDBC.DB.Begin()

	setUpUser(t, tx)

	tx.Rollback()
}

func Test_isExistProject(t *testing.T) {
	tx := TestDBC.DB.Begin()

	_, project := setUpUntilProject(t, tx)

	err := isExistProject(tx, project.Name, project.ID)
	if err != nil {
		t.Error(err)
	}

	tx.Rollback()
}

func TestGetProject(t *testing.T) {
	tx := TestDBC.DB.Begin()

	_, project := setUpUntilProject(t, tx)

	_, getProject, err := GetProject(tx, project.ID)
	if err != nil {
		t.Error(err)
	}

	if getProject.ID != project.ID {
		t.Error("取得したプロジェクトが正しくない。")
	}

	tx.Rollback()
}

func TestGetProjects(t *testing.T) {
	tx := TestDBC.DB.Begin()

	user, _ := setUpUntilProject(t, tx)

	_, _, err := CreateProject(tx, "test_second", TestProjectDescription, user.ID)
	if err != nil {
		t.Error(err)
	}

	_, projects, err := GetProjects(tx, user.ID)
	if err != nil {
		t.Error(err)
	}

	if len(projects) == 1 {
		t.Error("プロジェクトの数が正しく取得できていない。")
	}

	tx.Rollback()
}

func TestCreateProject(t *testing.T) {
	tx := TestDBC.DB.Begin()

	setUpUntilProject(t, tx)

	tx.Rollback()
}

func TestGetFolders(t *testing.T) {
	tx := TestDBC.DB.Begin()

	user, project, folder := setUpUntilFolder(t, tx)

	_, err := CreateFolder(tx, "second", user.ID, project.ID, &folder.ID)
	if err != nil {
		t.Error(err)
	}

	_, folders, err := GetFolders(tx, user.ID, project.ID)
	if err != nil {
		t.Error(err)
	}

	if len(folders) <= 1 {
		t.Error("プロジェクトに紐付けれているフォルダ数と取得したフォルダ数が異なる。")
	}
	tx.Rollback()
}

func TestCreateFolder(t *testing.T) {
	tx := TestDBC.DB.Begin()

	_, project, folder := setUpUntilFolder(t, tx)

	folder2, err := CreateFolder(tx, "second", user.ID, project.ID, &folder.ID)
	if err != nil {
		t.Error(err)
	}

	if *folder2.ParentFolderID != folder.ID {
		t.Error("作成されたフォルダが正しく紐付けれていない。")
	}

	tx.Rollback()
}

func TestDeleteFolder(t *testing.T) {
	tx := TestDBC.DB.Begin()

	user, project, folder := setUpUntilFolder(t, tx)

	err := DeleteFolder(tx, folder.ID)
	if err != nil {
		t.Error(err)
	}

	_, folders, err := GetFolders(tx, user.ID, project.ID)
	if err != nil {
		t.Error(err)
	}

	if len(folders) != 0 {
		t.Error("フォルダを削除できていない。")
	}

	tx.Rollback()
}

func TestGetFile(t *testing.T) {
	tx := TestDBC.DB.Begin()

	_, _, _, file := setUpUntilFile(t, tx)
	f, err := GetFile(tx, file.ID)
	if err != nil {
		t.Error(err)
	}

	if f.ID != file.ID {
		t.Error("取得したファイルが正しくない。")
	}

	tx.Rollback()
}

func TestGetFiles(t *testing.T) {
	tx := TestDBC.DB.Begin()

	user, project, folder, _ := setUpUntilFile(t, tx)
	_, err := CreateFile(tx, "second", user.ID, project.ID, &folder.ID)
	if err != nil {
		t.Error(err)
	}

	files, err := GetFiles(tx, user.ID, project.ID, folder.ID)
	if err != nil {
		t.Error(err)
	}

	if len(files) != 2 {
		t.Error("フォルダに紐付けらている全てのファイルを取得できていない。")
	}

	tx.Rollback()
}

func TestGetNoFolderFiles(t *testing.T) {
	tx := TestDBC.DB.Begin()

	user, project, folder := setUpUntilFolder(t, tx)
	CreateFile(tx, "no folder", user.ID, project.ID, nil)
	CreateFile(tx, "no folder2", user.ID, project.ID, nil)

	nff, err := GetNoFolderFiles(tx, user.ID, project.ID)
	if err != nil {
		t.Error(err)
	}

	files, err := GetFiles(tx, user.ID, project.ID, folder.ID)
	if err != nil {
		t.Error(err)
	}

	if len(nff) != 2 || len(files) != 0 {
		t.Error("ファイルがフォルダに紐付けられていないはずなのに紐付けられている。")
	}

	tx.Rollback()
}

func TestUpdateFile(t *testing.T) {
	tx := TestDBC.DB.Begin()

	_, _, _, file := setUpUntilFile(t, tx)

	beforeContent := file.Content

	err := UpdateFile(tx, file.ID, "update")
	if err != nil {
		t.Error(err)
	}

	f, err := GetFile(tx, file.ID)
	if err != nil {
		t.Error(err)
	}

	if beforeContent == f.Content {
		t.Error("ファイルのContentの変更が行われていない。")
	}

	tx.Rollback()
}

func TestCreateFile(t *testing.T) {
	tx := TestDBC.DB.Begin()

	setUpUntilFile(t, tx)

	tx.Rollback()
}

func TestDeleteFile(t *testing.T) {
	tx := TestDBC.DB.Begin()

	_, _, _, file := setUpUntilFile(t, tx)

	err := DeleteFile(tx, file.ID)
	if err != nil {
		t.Error(err)
	}

	f := &File{}
	tx = tx.Find(f, file.ID)

	if tx.RowsAffected != 0 {
		t.Error("ファイルが消去できていない。")
	}

	tx.Rollback()
}
