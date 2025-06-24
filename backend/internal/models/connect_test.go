package models

import (
	"fmt"
	"os"
	"testing"
)

// ↓NewDBConnector()の内部で"./.env"ファイルを探しにいくのでmodelsパッケージの配下に、
// テスト用の".env"ファイルを用意する必要がある。
const (
	testEnvPath = ".env"
)

var (
	testDSN          string = "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	dsnElmyArray     []any  = []any{"HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "SSL_MODE", "TZ"}
	testDsnElmyArray []any  = []any{"HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "SSL_MODE", "TZ"}
	TestDBC          DBConnector

	TestUserName     string = "test_user"
	TestUserEmail    string = "test@test"
	TestUserPassword string = "y6PEKW29"
	TestUserSalt     string = ""
	testUserData     *User  = &User{
		Username: TestUserName,
		Email:    TestUserEmail,
		Password: TestUserPassword,
		Salt:     TestUserSalt,
	}

	TestProjectName        string = "test_project"
	TestProjectDescription string = "this is the description"
	TestFolderName         string = "test_folder"
	TestFileName           string = "test_file"
)

func TestMain(m *testing.M) {
	fmt.Println("DB関連のテストに必要なデータの準備開始")
	TestDBC = *NewDBConnector(dsnElmyArray)

	tx := TestDBC.DB.Begin()
	if tx.Error != nil {
		fmt.Println(tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	_, testUser, err := CreateUser(tx, testUserData.Username, testUserData.Email, testUserData.Password, testUserData.Salt)
	if err != nil {
		fmt.Println("テストに必要なユーザーデータの作成に失敗しました。")
	}
	_, _, err = CreateProject(tx, TestProjectName, "", testUser.ID)
	if err != nil {
		fmt.Println("テストに必要なプロジェクトデータの作成に失敗しました。")
	}

	exitVal := m.Run()

	tx.Rollback()
	fmt.Println("DB関連のテスト終了")
	os.Exit(exitVal)
}

func TestCreateDSN(t *testing.T) {
	TestDBC.CreateDSN(testDsnElmyArray)
	if len(TestDBC.DSN) == len(testDSN) {
		t.Fatalf("%sファイルから値を取得できていない。", testEnvPath)
	}
}

func TestConnectDB(t *testing.T) {
	TestDBC.ConnectDB()
	if TestDBC.DB.Error != nil {
		t.Fatal("DB接続不可: ", TestDBC.DB.Error.Error())
	}

}
