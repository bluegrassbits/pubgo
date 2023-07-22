package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"pubgo/config"
	"pubgo/content"

	"github.com/gomarkdown/markdown"
)

// build non collection page
func buildNonCollectionPage(page config.Page) {
	log.Printf("Building page: %+v", page)
	pageFilename := cfg.ContentDir + "/" + page.Name + ".md"
	md, err := os.ReadFile(pageFilename)
	var entry content.Entry

	if err != nil {
		log.Println("Error reading markdown file:", err)
	}

	entry, md, _ = content.ParseEntry(md)

	var title string
	if page.Path == "/" {
		title = cfg.Site.Name + " ~ " + cfg.Site.Title
	} else {
		if entry.Title != "" {
			title = cfg.Site.Name + " ~ " + entry.Title
		} else {
			title = cfg.Site.Name + " ~ " + page.Name
		}
	}

	renderer, p := newCustomizedRender(entry.IncludeToc, cfg.Site.Theme.SyntaxHighlight)

	entry.Body = template.HTML(markdown.ToHTML(md, p, renderer))
	cont := content.Content{
		Site:        cfg.Site,
		Page:        page,
		RequestPath: page.Path,
		BasePath:    cfg.BaseURL,
		Mode:        cfg.Mode,
		Title:       title,
		Collection:  page.Collection,
		Entry:       entry,
	}

	primeDirectory(filepath.Join(cfg.OutputDir, page.Path))

	// file writer for index.html
	wr, err := os.Create(cfg.OutputDir + page.Path + "/index.html")
	if err != nil {
		log.Println("Error creating file:", err)
	}

	err = templates.ExecuteTemplate(wr, "indexHTML", cont)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func loadSingleEntry(page config.Page) {
	// Get markdown file
	data, err := os.ReadFile(filepath.Join(cfg.ContentDir, page.Name+".md"))

	// if file doesn't exist, create it with default content
	if os.IsNotExist(err) {
		data = []byte("# " + page.Name + "\n\n" + page.Name + " content goes here")
		err = os.WriteFile(filepath.Join(cfg.ContentDir, page.Name+".md"), data, 0644)
		if err != nil {
			log.Println("Error creating entry file:", err)
			panic(err)
		}
	}

	entry := createEntry(page, "", page.Name+".md", data)
	entries[page.Name] = append(entries[page.Name], entry)
}
