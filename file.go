package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// walkAndCopyFiles copies files from src to dest directory
func walkAndCopyFiles(src string, dest string) error {
	err := filepath.Walk(filepath.Join(src, "static"), func(path string, info os.FileInfo, err error) error {
		// sanitize content dir
		contentDir := strings.TrimPrefix(src, "./")

		// path removing content dir
		newPath := strings.Replace(path, contentDir, dest, 1)
		log.Printf("path: %s\n", path)
		log.Printf("newPath: %s\n", newPath)

		if info.IsDir() {
			// create directory in outputdir/static if it doesn't exist
			if info.Name() != "static" {
				err := os.MkdirAll(newPath, 0755)
				if err != nil {
					log.Println("Error creating static directory:", err)
					return err
				}
			}
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Println("Error reading file:", err)
				return err
			}

			err = os.WriteFile(newPath, data, 0644)
			if err != nil {
				log.Println("Error writing file:", err)
				return err
			}
		}

		return nil
	})

	return err
}
