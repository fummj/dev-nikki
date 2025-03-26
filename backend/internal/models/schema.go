package models

import (
	"gorm.io/gorm"

	"dev_nikki/internal/logger"
)

type User struct {
	ID       uint   `gorm:"primarykey;autoIncrement;not null"`
	Username string `gorm:"type:varchar(30);not null"`
	Email    string `gorm:"type:varchar(254);unique;not null"`
	Password string `gorm:"type:varchar(100)"`
	Salt     string `gorm:"type:varchar(16)"`
	IsActive bool   `gorm:"default:false"`
	gorm.Model
}

type Project struct {
	ID          uint   `gorm:"primarykey;autoIncrement;not null"`
	Name        string `gorm:"type:varchar(30);not null"`
	Description string `gorm:"type:varchar(100)"`
	UserID      uint   `gorm:"not null"`
	gorm.Model
}

type Folder struct {
	ID             uint   `gorm:"primarykey;autoIncrement;not null"`
	Name           string `gorm:"type:varchar(30);not null"`
	UserID         uint   `gorm:"not null"`
	ProjectID      uint   `gorm:"not null"`
	ParentFolderID uint
	gorm.Model
}

type File struct {
	ID        uint    `gorm:"primarykey;autoIncrement;not null"`
	Name      string  `gorm:"type:varchar(30);not null"`
	Content   *string `gorm:"type:text"`
	UserID    uint    `gorm:"not null"`
	ProjectID uint    `gorm:"not null"`
	FolderID  uint
	gorm.Model
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
