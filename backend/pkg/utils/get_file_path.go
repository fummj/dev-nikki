package utils

import (
	"fmt"
	"os"
)

func SearchFileFindParentDir(dir, name string) string {
	d := ""
	for {
		if f, err := os.Stat(fmt.Sprint(dir, "/", d, name)); os.IsNotExist(err) {
			d = fmt.Sprint(d, "../")
			continue
		} else {
			if !f.IsDir() {
				return fmt.Sprint(dir, "/", d, name)
			}
		}
	}
}
