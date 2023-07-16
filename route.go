package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pubgo/config"
	"pubgo/content"
	"strings"

	"github.com/gomarkdown/markdown"
)

var fourOhFour = `
<h2>404: Page not found.</h2>
<pre>
 _________________
< You Done Goofed >
 -----------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
</pre>
`

// setup main router
func setupRouter() {
	// a handler to process the request path and map it to a page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		path := r.URL.Path
		// generate filepath from path
		route, err := parseRoute(path)

		log.Printf("Route: %s, Error: %s\n", route, err)

		if err != nil {
			handleNotFoundError(w, r)
			return
		}

		// if route is not a directory
		if !isDir(route) {
			if isSinglePage(route) {
				renderSinglePage(w, path, route)
				return
			} else {
				renderEntryPage(w, r, route)
				return
			}
		} else {
			renderEntriesPage(w, r, route)
			//fmt.Fprintln(w, "Collection handler")
			return
		}
	})
}

// isSinglePage determines if the route is a single page or an entry in a collection
// it does this by checking if there is a subdirectory in the path relative to the content directory

func isSinglePage(path string) bool {
	contentDir := cfg.ContentDir

	// if cfg.ContentDir has "./" prefix then remove it
	if strings.HasPrefix(contentDir, "./") {
		contentDir = strings.TrimPrefix(contentDir, "./")
	}

	// remove contentDir from path
	path = strings.TrimPrefix(path, contentDir)
	fileParts := strings.Split(path, string(os.PathSeparator))

	log.Printf("Path: %s, FileParts: %s\n", path, fileParts)
	if len(fileParts) > 2 {
		if isDir(filepath.Join(cfg.ContentDir, fileParts[0])) {
			return false
		}
		return false
	}
	return true
}

// isDir returns true if the path is a directory
func isDir(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.IsDir()
	}
	return false
}

// isRootPath returns true if the path is the root path ("/" or "/index.html")
func isRootPath(path string) bool {
	return path == "/" || path == "" || path == "/index.html" || path == "/index"
}

func parseRoute(path string) (string, error) {
	// if the path is "/" or "/index.html" then use the home page
	if isRootPath(path) {
		if _, err := os.Stat(filepath.Join(cfg.ContentDir, "home.md")); err == nil {
			return filepath.Join(cfg.ContentDir, "home.md"), nil
		} else {
			return "", err
		}
	}

	// if the path doesn't have an extension then it's a page
	if filepath.Ext(path) == "" {
		// if path + ".md" exists then use it
		if _, err := os.Stat(filepath.Join(cfg.ContentDir, path+".md")); err == nil {
			return filepath.Join(cfg.ContentDir, path+".md"), nil
		} else if _, err := os.Stat(filepath.Join(cfg.ContentDir, path)); err == nil {
			return filepath.Join(cfg.ContentDir, path), nil
		}
	}

	// if the path has an extension then it's a file
	if filepath.Ext(path) != "" {
		// if html lookup md file
		if filepath.Ext(path) == ".html" {
			path = strings.TrimSuffix(path, filepath.Ext(path))
			if _, err := os.Stat(filepath.Join(cfg.ContentDir, path+".md")); err == nil {
				return filepath.Join(cfg.ContentDir, path+".md"), nil
			}
		} else if _, err := os.Stat(filepath.Join(cfg.ContentDir, path)); err == nil {
			return filepath.Join(cfg.ContentDir, path), nil
		}
	}

	return "", fmt.Errorf("Route not found")
}

// handleNotFoundError handles the request for a non-existing route
func handleNotFoundError(w http.ResponseWriter, r *http.Request) {
	cont := content.Content{
		Site:        cfg.Site,
		RequestPath: r.URL.Path,
		BasePath:    cfg.BaseURL,
		Mode:        cfg.Mode,
		Title:       cfg.Site.Name + " ~ 404",
		Entry: content.Entry{
			Title: "404",
			Body:  template.HTML(fourOhFour),
		},
	}
	w.WriteHeader(http.StatusNotFound)
	err := templates.ExecuteTemplate(w, "indexHTML", cont)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderSinglePage renders a single page (Markdown or HTML) to the response writer
func renderSinglePage(w http.ResponseWriter, requestPath string, filePath string) {
	pages := cfg.Site.Pages
	var page config.Page
	for _, p := range pages {
		if p.Path == requestPath {
			page = p
		}
	}

	log.Printf("Page: %+v\n", page)
	md, err := os.ReadFile(filePath)
	entry := content.Entry{}

	if err != nil {
		http.Error(w, "Error reading markdown file: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error reading markdown file:", err)
		return
	}

	entry, md, _ = content.ParseEntry(md)
	var title string

	if requestPath == "/" {
		title = cfg.Site.Name + " ~ " + cfg.Site.Title
	} else {
		if entry.Title != "" {
			title = cfg.Site.Name + " ~ " + entry.Title
		} else {
			title = cfg.Site.Name + " ~ " + page.Name
		}
	}

	renderer, p := newCustomizedRender(page.Toc)
	entry.Body = template.HTML(markdown.ToHTML(md, p, renderer))
	cont := content.Content{
		Site:        cfg.Site,
		Page:        page,
		RequestPath: requestPath,
		Mode:        cfg.Mode,
		Title:       title,
		Collection:  page.Collection,
		Entry:       entry,
	}

	err = templates.ExecuteTemplate(w, "indexHTML", cont)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderEntryPage renders an entry page (Markdown or HTML) to the response writer
func renderEntryPage(w http.ResponseWriter, r *http.Request, filePath string) {
	collection := filepath.Base(filepath.Dir(filePath))
	collections := cfg.Site.Pages
	var page config.Page
	for _, c := range collections {
		if c.Name == collection {
			page = c
		}
	}

	md, err := os.ReadFile(filePath)
	entry := content.Entry{}

	if err != nil {
		http.Error(w, "Error reading markdown file: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error reading markdown file:", err)
		return
	}

	entry, md, _ = content.ParseEntry(md)
	var title string

	if entry.Title != "" {
		title = cfg.Site.Name + " ~ " + entry.Title
	} else {
		title = cfg.Site.Name + " ~ " + page.Name
	}

	renderer, p := newCustomizedRender(page.Toc)
	entry.Body = template.HTML(markdown.ToHTML(md, p, renderer))
	cont := content.Content{
		Site:        cfg.Site,
		RequestPath: r.URL.Path,
		Mode:        cfg.Mode,
		Title:       title,
		Page:        page,
		Entry:       entry,
	}
	if r.Header.Get("HX-Request") == "true" {
		cont.Title = cfg.Site.Name + " ~ " + r.URL.Path

		err := templates.ExecuteTemplate(w, "entryHTML", cont)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		cont.Title = cfg.Site.Name + " - " + entry.Title

		err := templates.ExecuteTemplate(w, "indexHTML", cont)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}

// renderEntriesPage renders an entries page (Markdown or HTML) to the response writer
func renderEntriesPage(w http.ResponseWriter, r *http.Request, filePath string) {
	pages := cfg.Site.Pages

	var page config.Page
	for _, p := range pages {
		if p.Path == r.URL.Path {
			page = p
		}
	}

	files, _ := ioutil.ReadDir(filePath)
	var ents []content.Entry
	for _, f := range files {
		log.Printf("File: %s/%s\n", filePath, f.Name())
		filename := f.Name()

		if filepath.Ext(f.Name()) == ".md" {
			data, err := os.ReadFile(filepath.Join(filePath, f.Name()))

			if err != nil {
				http.Error(w, "Error reading markdown file: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error reading markdown file:", err)
				return
			}
			entry := createEntry(page, page.Name, filename, data)
			ents = append(ents, entry)
		}
	}

	cont := createContent(page, ents)
	err := templates.ExecuteTemplate(w, "indexHTML", cont)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
