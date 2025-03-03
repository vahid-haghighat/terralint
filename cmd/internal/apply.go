package internal

import (
	"io/fs"
	"os"
	"path/filepath"
)

func Apply(filePath string) error {
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return applyRulesToDirectory(filePath)
	}
	return applyRulesToFile(filePath)
}

func applyRulesToFile(filePath string) error {
	extension := filepath.Ext(filePath)

	if extension != ".tf" && extension != ".tfvars" {
		return nil
	}
	formattedBytes, err := getFormattedContent(filePath)

	if err != nil {
		return err
	}

	return os.WriteFile(filePath, formattedBytes, 0666)
}

func applyRulesToDirectory(directoryPath string) error {
	err := filepath.WalkDir(directoryPath, func(path string, d fs.DirEntry, err error) error {
		if path == directoryPath {
			return nil
		}
		stat, _ := os.Stat(path)
		if stat.IsDir() {
			return applyRulesToDirectory(path)
		}
		return applyRulesToFile(path)
	})
	if err != nil {
		return err
	}
	return nil
}
