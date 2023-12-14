package cmd

import (
	"fmt"
	"os"
)

func ListFiles(dirname string) ([]string, error) {
	direntries, err := os.ReadDir(dirname)
	
	if err != nil {
		return []string{}, fmt.Errorf("Unable to read directory %v: %w", dirname, err)
	}
	
	result := make([]string, len(direntries))
	
	for index, entry := range direntries {
		result[index] = entry.Name()
	}
	return result, nil
}

func Execute(dirname string) error {
	
	fileinfo, err := os.Stat(dirname)
	if err != nil {
		return err
	}

	if !fileinfo.IsDir() {
		fmt.Println(dirname)
		return nil
	}
	
	files, err := ListFiles(dirname)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}
