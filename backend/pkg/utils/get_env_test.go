package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var envPath = ".env.test"

func TestGetEnv(t *testing.T) {

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	envPath := filepath.Join(currentDir, ".env.test")

	testEnvData := `
		TEST_USER="testdesu"
		TEST_PASSWORD="DWQ8VwlQDodU"
		TEST_DB_NAME="test_db"
	`

	err := os.WriteFile(envPath, []byte(testEnvData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(envPath)

	expected := map[string]string{
		"TEST_USER":     "testdesu",
		"TEST_PASSWORD": "DWQ8VwlQDodU",
		"TEST_DB_NAME":  "test_db",
	}

	result := GetEnv(envPath)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}
