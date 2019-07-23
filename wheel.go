package main

import (
	"os"
	"path/filepath"
	"strconv"
)

func dirSize(path string) (size int64, err error) {
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func rmin(size int64) string {
	return strconv.FormatInt(1, 10)
}

func rmax(size int64) string {
	return strconv.FormatInt(size/(speed*seconds)+1, 10)
}
