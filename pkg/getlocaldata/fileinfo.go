package getlocaldata

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FindLastModifiedFile - find all files with search term and get last modified file
func FindLastModifiedFile(dir, searchTerm string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		log.Println("Error to get list files", err)
		return "", err
	}

	var filteredFiles []string
	for _, file := range files {
		if strings.Contains(file, searchTerm) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	if len(filteredFiles) == 0 {
		fmt.Println("Files didn't find")
		return "", err
	}

	var lastModifiedFile string
	var lastModifiedTime time.Time
	for _, file := range filteredFiles {
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Println("Error get info from file", err)
			return "", err
		}

		modTime := fileInfo.ModTime()
		if modTime.After(lastModifiedTime) {
			lastModifiedTime = modTime
			lastModifiedFile = file
		}
	}

	return lastModifiedFile, nil
}
