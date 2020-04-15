package main

import (
	"os"
	"path/filepath"
	"strings"
)

const ova = ".ova"

func listCachedImages() ([]string, error) {
	baseDir := getViperHomeDir() + "/boxes/"

	images := make([]string, 0)

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ova) {
			name := strings.ReplaceAll(path, baseDir, "")
			images = append(images, strings.ReplaceAll(name, ova, ""))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return images, nil
}
