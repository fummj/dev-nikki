package utils

import (
	"os"
)

var envMap map[string]string = make(map[string]string, 0)

// セットされている環境変数を取得する。
func GetEnv(s []any) map[string]string {
	m := map[string]string{}

	for _, v := range s {
		key, _ := v.(string)
		m[key] = os.Getenv(key)
	}

	return m
}
