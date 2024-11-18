package sys

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindPath(fileName string, paths []string) (fileLoc string) {
	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory " + path)
			continue
		}
		for _, file := range files {
			if file.Name() == fileName {
				fileLoc = path + "/" + fileName
				return
			}
		}
	}
	fmt.Printf("File %s not found in any of the paths\n", fileName)
	return
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func GetFileList(root string) (paths []string, err error) {
	paths = make([]string, 0)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsDir() {
			return nil
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})
	return paths, err
}
