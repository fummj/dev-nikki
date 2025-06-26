package utils

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetEnv(t *testing.T) {

	testEnvData := map[string]string{
		"TEST_USER":     "testdesu",
		"TEST_PASSWORD": "DWQ8VwlQDodU",
		"TEST_DB_NAME":  "test_db",
	}

	// 環境変数を設定
	for k, v := range testEnvData {
		os.Setenv(k, v)
	}

	// 設定した環境変数を削除
	defer func() {
		for k := range testEnvData {
			os.Unsetenv(k)
		}
	}()

	expected := map[string]string{
		"TEST_USER":     "testdesu",
		"TEST_PASSWORD": "DWQ8VwlQDodU",
		"TEST_DB_NAME":  "test_db",
	}

	keys := []any{"TEST_USER", "TEST_PASSWORD", "TEST_DB_NAME"}
	result := GetEnv(keys)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}
