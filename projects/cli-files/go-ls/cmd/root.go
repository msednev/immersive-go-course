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

func Execute() {
	files, err := ListFiles(".")

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	for _, file := range files {
		fmt.Println(file)
	}
}
