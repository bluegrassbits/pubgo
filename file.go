package main

import (
	"log"
	"os"
	"path/filepath"
	"pubgo/config"
	"strings"
)

func primeDirectory(dir string) {
	// If directory doesn't exist, create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Println("Error creating content directory:", err)
			panic(err)
		}
	}
}

func primePageEntry(page config.Page) {
	// Get markdown file
	data, err := os.ReadFile(filepath.Join(cfg.ContentDir, page.Name+".md"))

	// if file doesn't exist, create it with default content
	if os.IsNotExist(err) {
		data = []byte("### " + page.Name + ".md\n\n" + page.Name + " content goes here. edit this file to change the page content.")
		err = os.WriteFile(filepath.Join(cfg.ContentDir, page.Name+".md"), data, 0644)
		if err != nil {
			log.Println("Error creating entry file:", err)
			panic(err)
		}
	}
}

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
