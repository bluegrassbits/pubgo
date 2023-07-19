package main

import (
	"html/template"
	"log"
	"os"

	"pubgo/config"
	"pubgo/content"

	"github.com/gomarkdown/markdown"
)

// buildCollectionPage builds a collection page.
func buildCollectionPage(page config.Page) {
	log.Printf("Building collection page: %s", page.Name)

	ents := entries[page.Name]
	cont := createContent(page, ents)

	createOutputDirectory(page)

	wr, err := os.Create(cfg.OutputDir + page.Path + "/index.html")
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}

	err = templates.ExecuteTemplate(wr, "indexHTML", cont)
	if err != nil {
		log.Println("Error executing template:", err)
	}

	buildEntryPages(page)
}

// buildEntryPages builds the pages for each entry in a collection.
func buildEntryPages(page config.Page) {
	log.Printf("Building entry pages for collection: %s", page.Name)

	ents := entries[page.Name]
	log.Printf("Entries: %+v", ents)

	for _, entry := range ents {

		entryBody := loadEntryBody(page, entry) // Updated here
		entry.Body = template.HTML(entryBody)

		cont := content.Content{
			Site:        cfg.Site,
			Page:        page,
			RequestPath: page.Path,
			BasePath:    cfg.BaseURL,
			Mode:        cfg.Mode,
			Title:       cfg.Site.Name + " ~ " + page.Name,
			Collection:  page.Collection,
			Entry:       entry,
		}

		createOutputDirectory(page)

		wr, err := os.Create(cfg.OutputDir + page.Path + "/" + entry.StaticFileName())
		if err != nil {
			log.Println("Error creating file:", err)
			return
		}

		err = templates.ExecuteTemplate(wr, "indexHTML", cont)
		if err != nil {
			log.Println("Error executing template:", err)
		}
	}
}

// createContent creates a Content struct for a collection page.
func createContent(page config.Page, ents []content.Entry) content.Content {
	if len(ents) > 0 {
		return content.Content{
			Site:        cfg.Site,
			Page:        page,
			RequestPath: page.Path,
			Mode:        cfg.Mode,
			BasePath:    cfg.BaseURL,
			Title:       cfg.Site.Name + " ~ " + page.Name,
			Collection:  page.Collection,
			Entries:     ents,
		}
	}

	return content.Content{
		Site:        cfg.Site,
		Page:        page,
		RequestPath: page.Path,
		Mode:        cfg.Mode,
		BasePath:    cfg.BaseURL,
		Title:       cfg.Site.Name + " ~ " + page.Name,
		Collection:  page.Collection,
		Entry: content.Entry{
			Body: template.HTML("<b>No entries found</b><p>Please create a new entry in this page.</p>"),
		},
	}
}

// findEntryByFilename finds an entry by filename in the collection.
func findEntryByFilename(ents []content.Entry, filename string) content.Entry {
	for _, entry := range ents {
		if entry.FileName == filename {
			return entry
		}
	}

	return content.Entry{}
}

// loadEntryBody loads the body of an entry.
func loadEntryBody(page config.Page, entry content.Entry) string {
	entryFilename := cfg.ContentDir + "/" + page.Name + "/" + entry.FileName // Updated here
	md, err := os.ReadFile(entryFilename)
	if err != nil {
		log.Println(err)
		return ""
	}

	_, md, _ = content.ParseEntry(md)
	renderer, p := newCustomizedRender(entry.IncludeToc, cfg.Site.Theme.SyntaxHighlight)

	mdString := string(md)
	html := markdown.ToHTML([]byte(mdString), p, renderer)

	return string(html)
}

// createOutputDirectory creates the output directory for a page if it doesn't exist.
func createOutputDirectory(page config.Page) {
	if _, err := os.Stat(cfg.OutputDir + page.Path); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.OutputDir+page.Path, 0755)
		if err != nil {
			log.Println("Error creating output directory:", err)
			panic(err)
		}
	}
}
