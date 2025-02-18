package utils

import (
	"log"

	"github.com/joho/godotenv"
)

var envMap map[string]string = make(map[string]string, 0)

// 環境変数入りのmapを取得。引数pは隠しファイルのpath
func GetEnv(p string) map[string]string {
	godotenv.Load(p)
	envMap, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	return envMap
}
