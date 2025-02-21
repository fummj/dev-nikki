package models

import (
	"gorm.io/gorm"

	"dev_nikki/internal/logger"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(30)"`
	Email    string `gorm:"type:varchar(254);unique"`
	Password string `gorm:"type:varchar(100)"`
	Salt     string `gorm:"type:varchar(16)"`
	IsActive bool   `gorm:"default:false"`
}

type Project struct {
	gorm.Model
	Name        string `gorm:"type:varchar(30)"`
	Description string `gorm:"type:varchar(100)"`
	UserID      uint
}

type Folder struct {
	gorm.Model
	Name      string `gorm:"type:varchar(30)"`
	UserID    uint
	ProjectID uint
}

type File struct {
	gorm.Model
	Name      string  `gorm:"type:varchar(30)"`
	Content   *string `gorm:"type:text"`
	UserID    uint
	ProjectID uint
	FolderID  uint
}

func IsExistTable(db *gorm.DB) bool {
	if !db.Migrator().HasTable(&User{}) {
		return false
	}
	return true
}

// テーブルの有無確認して、存在しない場合は作成
func FirstMigration(db *gorm.DB) {
	if IsExistTable(db) {
		logger.Slog.Debug("the table already exists", "method", "FirstMigration")
		return
	}

	err := db.AutoMigrate(
		&User{},
		&Project{},
		&Folder{},
		&File{},
	)
	if err != nil {
		logger.Slog.Error("Failed: "+err.Error(), "method", "FirstMigration")
	}
	logger.Slog.Info("all tables have been initialized", "method", "FirstMigration")
}

// 全テーブル削除
func AllDropTables(db *gorm.DB) {
	err := db.Migrator().DropTable(
		&User{},
		&Project{},
		&Folder{},
		&File{},
	)
	if err != nil {
		logger.Slog.Error("Failed: "+err.Error(), "method", "AllDropTables")
		return
	}
	logger.Slog.Info("all droped", "method", "AllDropTables")
}
