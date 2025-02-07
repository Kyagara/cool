package utils

import (
	"os"
)

func GetFileSize(path string) int {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return int(stat.Size())
}
