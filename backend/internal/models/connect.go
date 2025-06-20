package models

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dev_nikki/internal/logger"
	"dev_nikki/pkg/utils"
)

const (
	EnvPath = ".env"
)

var (
	dsn          string = "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	DsnElmyArray []any  = []any{"HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "SSL_MODE", "TZ"}
	DBC                 = NewDBConnector(DsnElmyArray)
)

func init() {
	FirstMigration(DBC.DB)
	// AllDropTables(DBC.DB)
}

type dbConnector interface {
	CreateDSN([]any)
	ConnectDB()
}

type DBConnector struct {
	DSN string
	DB  *gorm.DB
}

// DSN作成
func (c *DBConnector) CreateDSN(s []any) {
	m := utils.GetEnv(EnvPath)
	for i := 0; i < len(s); i++ {
		if v, ok := s[i].(string); ok {
			s[i] = m[v]
		} else {
			logger.Slog.Error("Failed: dsnElmyArray elements contain not string")
			return
		}
	}
	c.DSN = fmt.Sprintf(dsn, s...)
}

// DB接続
func (c *DBConnector) ConnectDB() {
	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	if err != nil {
		logger.Slog.Error("Failed: " + err.Error())
		return
	}
	c.DB = db
}

// インスタンスを返す。引数sは接続情報のkeyを順番通りに入れたslice
func NewDBConnector(s []any) *DBConnector {
	con := DBConnector{}
	con.CreateDSN(s)
	con.ConnectDB()
	return &con
}
