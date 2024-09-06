package sys

import (
	"fmt"
	"log"
	"os"
)

func FilePath(fileName string, paths []string) (fileLoc string) {
	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory " + path)
			continue
		}
		for _, file := range files {
			if file.Name() == fileName {
				fileLoc = path + "/" + fileName
				log.Println(fileName+" found at: ", fileLoc)
				return
			}
		}
	}
	fmt.Printf("File %s not found in any of the paths\n", fileName)
	return
}
