package utils

import (
	"fmt"
	"os"
)

func GetFilePath(filename string) string {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename
	}

	if _, err := os.Stat(fmt.Sprint("./../../", filename)); !os.IsNotExist(err) {
		return fmt.Sprint("./../../", filename)
	}

	return ""
}
