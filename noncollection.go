package main

import (
	"html/template"
	"log"
	"net/http"
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

	renderer := newCustomizedRender(page.Toc)

	entry.Body = template.HTML(markdown.ToHTML(md, nil, renderer))
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

func setupNonCollectionPageHandler(page config.Page) {
	log.Println("Setting up handler for page:", page.Path)

	pageFilename := cfg.ContentDir + "/" + page.Name + ".md"

	http.HandleFunc(page.Path, func(wr http.ResponseWriter, req *http.Request) {
		handleNonCollectionPageRequest(wr, req, page, pageFilename)
	})
}

func handleNonCollectionPageRequest(wr http.ResponseWriter, req *http.Request, page config.Page, pageFilename string) {

	if page.Path == "/" && req.URL.Path != "/" {
		cont := content.Content{
			Site:        cfg.Site,
			Page:        page,
			RequestPath: req.URL.Path,
			BasePath:    cfg.BaseURL,
			Mode:        cfg.Mode,
			Title:       cfg.Site.Name + " ~ 404",
			Entry: content.Entry{
				Title: "404",
				Body:  template.HTML(fourOhFour),
			},
		}

		// set status code to 404
		wr.WriteHeader(http.StatusNotFound)

		err := templates.ExecuteTemplate(wr, "indexHTML", cont)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	md, err := os.ReadFile(pageFilename)
	entry := content.Entry{}

	if err != nil {
		http.Error(wr, "Error reading markdown file: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error reading markdown file:", err)
		return
	}

	entry, md, _ = content.ParseEntry(md)

	var title string
	if req.URL.Path == "/" {
		title = cfg.Site.Name + " ~ " + cfg.Site.Title
	} else {
		if entry.Title != "" {
			title = cfg.Site.Name + " ~ " + entry.Title
		} else {
			title = cfg.Site.Name + " ~ " + page.Name
		}
	}

	renderer := newCustomizedRender(page.Toc)

	entry.Body = template.HTML(markdown.ToHTML(md, nil, renderer))
	cont := content.Content{
		Site:        cfg.Site,
		Page:        page,
		RequestPath: req.URL.Path,
		Mode:        cfg.Mode,
		Title:       title,
		Collection:  page.Collection,
		Entry:       entry,
	}

	err = templates.ExecuteTemplate(wr, "indexHTML", cont)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
	}
}
