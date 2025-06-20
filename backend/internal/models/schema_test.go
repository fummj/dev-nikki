package models

import (
	"testing"
)

// func TestIsExistTable(t *testing.T) {
// 	result := IsExistTable(TestDBC.DB)
// 	if !result {
// 		t.Error("DB内に必要なテーブルが作成されていない。")
// 	}
// }

// func TestFirstMigration(t *testing.T) {
// 	FirstMigration(TestDBC.DB)
// 	if TestDBC.DB.Error != nil {
// 		t.Error("テーブルの初期化に失敗しました。")
// 	}
// }

func TestAllDropTable(t *testing.T) {
	// trunsaction + rollbackが必要。
	// trunsactionがデフォルトで作成, 更新, 削除に対して有効になっている。
	// 自分でトランザクションやロールバックをする場合は、
	// TestDBC.DBで行うのではなくトランザクションオブジェクトを使用する必要がある。

	// tx := TestDBC.DB.Begin()
	// tx.SavePoint("sp1")
	//
	// AllDropTables(tx)
	// if tx.Error != nil {
	// 	t.Error("DB内の全てのテーブルの削除が失敗しました。")
	// }
	//
	// tx.RollbackTo("sp1")
	//
	// tx.Commit()
}
