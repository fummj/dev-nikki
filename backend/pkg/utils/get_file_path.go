package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func SearchFileFindParentDir(dir, name string) string {
	for {
		if f, err := os.Stat(fmt.Sprint(dir, "/", name)); os.IsNotExist(err) {
			dir = filepath.Dir(dir)
			continue
		} else {
			if !f.IsDir() {
				return fmt.Sprint(dir, "/", name)
			}
		}
	}
}
