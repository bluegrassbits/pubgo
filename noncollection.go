package main

import (
	"html/template"
	"log"
	"os"

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

	renderer, p := newCustomizedRender(entry.IncludeToc)

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

	// if output dir/path doesn't exist, create it
	if _, err := os.Stat(cfg.OutputDir + page.Path); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.OutputDir+page.Path, 0755)
		if err != nil {
			log.Println("Error creating output directory:", err)
			panic(err)
		}
	}

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
