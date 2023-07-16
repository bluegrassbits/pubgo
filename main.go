package main

import (
	"embed"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pubgo/config"
	"pubgo/content"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS
var entries = make(map[string][]content.Entry)
var cfg = config.NewConfig()
var templates *template.Template

var runMode *string
var outputDir *string

func init() {
	// Load config from YAML file
	configFile := flag.String("config", "config.yaml", "Path to config file")
	runMode = flag.String("mode", "serve", "Run mode: <serve> or <build> static site")
	outputDir = flag.String("out", "./out", "Output directory for static site")
	contentDir := flag.String("content_dir", "./website", "Content directory")

	flag.Parse()
	cfg.ContentDir = *contentDir
	cfg.OutputDir = *outputDir
	cfg.Mode = *runMode
	config.LoadConfig(*configFile, &cfg)

	// If content directory doesn't exist, create it
	if _, err := os.Stat(cfg.ContentDir); os.IsNotExist(err) {
		err = os.Mkdir(cfg.ContentDir, 0755)
		if err != nil {
			log.Println("Error creating content directory:", err)
			panic(err)
		}
	}

	// if runmode is build and output directory doesn't exist, create it
	if cfg.Mode == "build" {
		if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
			err = os.Mkdir(cfg.OutputDir, 0755)
			if err != nil {
				log.Println("Error creating output directory:", err)
				panic(err)
			}
		}
	}

	// Load templates
	var err error
	templates, err = template.New("").ParseFS(templateFiles, "templates/*.tmpl")

	if err != nil {
		// If error loading templates, log it
		log.Println("Error loading templates:", err)
	}

	// Load custom templates from ContentDir/templates to override default templates
	customTemplatesDir := filepath.Join(cfg.ContentDir, "templates")

	// If custom templates directory doesnt exist, create it
	if _, err := os.Stat(customTemplatesDir); os.IsNotExist(err) {
		err = os.Mkdir(customTemplatesDir, 0755)
		if err != nil {
			log.Println("Error creating custom templates directory:", err)
			panic(err)
		}
	}

	if _, err := os.Stat(customTemplatesDir); err == nil {
		// Log if custom templates directory exists
		log.Println("Loading custom templates from", customTemplatesDir)

		// Print *.html files in custom templates directory
		files, err := ioutil.ReadDir(customTemplatesDir)
		if err != nil {
			log.Println("Error reading custom templates directory:", err)
		}

		var count int
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".tmpl") {
				log.Println("Found custom template:", file.Name())
				count++
			}
		}

		if count == 0 {
			log.Println("No custom templates found")
		} else {
			log.Println("Found", count, "custom templates")
			templates, err = templates.New("").ParseGlob(filepath.Join(customTemplatesDir, "*.tmpl"))
			if err != nil {
				log.Println("Error loading custom templates:", err)
			}
		}
	}

	// Load entries
	for _, page := range cfg.Site.Pages {
		if page.Collection {
			loadCollectionEntries(page)
		} else {
			loadSingleEntry(page)
		}
	}

	// Print loaded entries
	printEntries()
}

func loadCollectionEntries(page config.Page) {
	// Get list of entry markdown files
	files, err := ioutil.ReadDir(filepath.Join(cfg.ContentDir, page.Name))
	// If directory doesn't exist, create it
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(cfg.ContentDir, page.Name), 0755)
		if err != nil {
			log.Println("Error creating page directory:", err)
			panic(err)
		}
		files = []os.FileInfo{}
	} else if err != nil {
		log.Println("Error reading page directory:", err)
		panic(err)
	}

	log.Println("Loading entries...")

	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".md") {
			data, err := os.ReadFile(filepath.Join(cfg.ContentDir, page.Name, filename))
			if err != nil {
				log.Println("Error reading entry file:", err)
				panic(err)
			}

			entry := createEntry(page, page.Name, filename, data)
			entries[page.Name] = append(entries[page.Name], entry)
		}
	}

	log.Println("Found", len(files), "entries")
	log.Println("Done loading entries")
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

func createEntry(page config.Page, subDir, filename string, data []byte) content.Entry {

	var filePath string
	if subDir != "" {
		filePath = filepath.Join(cfg.ContentDir, page.Name, filename)
	} else {
		filePath = filepath.Join(cfg.ContentDir, filename)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading entry file:", err)
		panic(err)
	}

	entry, _, err := content.ParseEntry(data)
	if err != nil {
		log.Println("Error parsing entry:", err)
	}

	// filename change .md to .html
	entry.FileName = strings.Replace(filename, ".md", ".html", 1)
	entry.Page = page.Name

	return entry
}

func printEntries() {
	for _, entryList := range entries {
		for _, entry := range entryList {
			log.Println(entry)
		}
	}
}

// log request
func logRequest(req *http.Request) {
	log.Printf("%s\t%s\t%s\t%s", req.Method, req.URL.Path, req.RemoteAddr, req.UserAgent())
}

// loadEntries
func loadEntries() {
	// Clear entries
	entries = make(map[string][]content.Entry, 0)

	// Load entries
	for _, page := range cfg.Site.Pages {
		if page.Collection {
			loadCollectionEntries(page)
		} else {
			loadSingleEntry(page)
		}
	}

	// Print loaded entries
	printEntries()
}

func main() {
	if cfg.Mode == "build" {

		// using os.Read and os.Write copy files from contentdir/static/ to outputdir/static/
		err := walkAndCopyFiles(cfg.ContentDir, cfg.OutputDir)

		if err != nil {
			log.Println("Error walking static directory:", err)
			panic(err)
		}

		// make outputdir/css if it doesn't exist
		if _, err := os.Stat(filepath.Join(cfg.OutputDir, "css")); os.IsNotExist(err) {
			err = os.Mkdir(filepath.Join(cfg.OutputDir, "css"), 0755)
			if err != nil {
				log.Println("Error creating css directory:", err)
				panic(err)
			}
		}

		// render css from template and write to outputdir/css/style.css
		wr, err := os.Create(cfg.OutputDir + "/css" + "/style.css")
		if err != nil {
			log.Println("Error creating file:", err)
		}

		err = templates.ExecuteTemplate(wr, "styleCSS", cfg.Site.Theme)
		if err != nil {
			log.Println("Error executing template:", err)
		}

		buildPages()
	}

	if cfg.Mode == "serve" {
		serveStaticFiles()
		serveCSSTemplate()
		setupRouter()

		// Start web server
		log.Println("Starting web server on port", cfg.Port)
		err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
		if err != nil {
			log.Fatal("Web server error:", err)
		}
	}

}

// buildPages builds the static pages
func buildPages() {
	nonCollectionPages := make(config.Pages)
	collectionPages := make(config.Pages)

	for _, page := range cfg.Site.Pages {
		if page.Collection {
			collectionPages[page.Name] = page
		} else {
			nonCollectionPages[page.Name] = page
		}
	}

	// Set up handlers for non-collection pages
	for _, page := range nonCollectionPages {
		buildNonCollectionPage(page)
	}

	// Set up handlers for collection pages
	for _, page := range collectionPages {
		buildCollectionPage(page)
	}
}
