package utils

import (
	"log"

	"github.com/joho/godotenv"
)

var envMap map[string]string = make(map[string]string, 0)

// 環境変数入りのmapを取得。引数pは隠しファイルのpath
func GetEnv(p string) map[string]string {
	err := godotenv.Load(p)
	if err != nil {
		log.Fatal(err)
	}

	envMap, err := godotenv.Read(p)
	if err != nil {
		log.Fatal(err)
	}
	return envMap
}
