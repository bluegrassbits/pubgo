package main

import (
	"crypto/subtle"
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

func init() {
	// Load config from YAML file
	configFile := flag.String("config", "config.yaml", "Path to config file")
	runMode := flag.String("mode", "serve", "Run mode: <serve> or <build> static site")
	outputDir := flag.String("out", "./out", "Output directory for static site")
	contentDir := flag.String("content_dir", "./website", "Content directory")
	adminUser := flag.String("admin_user", "", "Admin username")
	adminPass := flag.String("admin_pass", "", "Admin password")

	flag.Parse()
	cfg.ContentDir = *contentDir
	cfg.OutputDir = *outputDir
	cfg.Mode = *runMode
	cfg.AdminUser = *adminUser
	cfg.AdminPass = *adminPass

	config.LoadConfig(*configFile, &cfg)
	primeDirectory(cfg.ContentDir)

	for _, page := range cfg.Site.Pages {
		if page.Collection {
			primeDirectory(filepath.Join(cfg.ContentDir, page.Name))
		} else {
			primePageEntry(page)
		}
	}

	// if runmode is build and output directory doesn't exist, create it
	if cfg.Mode == "build" {
		primeDirectory(cfg.OutputDir)

		// Load entries
		for _, page := range cfg.Site.Pages {
			if page.Collection {
				loadCollectionEntries(page)
			} else {
				loadSingleEntry(page)
			}
		}

		printEntries()
	}

	// Load default templates
	var err error
	templates, err = template.New("").ParseFS(templateFiles, "templates/*.tmpl")

	if err != nil {
		// If error loading templates, log it
		log.Println("Error loading templates:", err)
	}

	// Load custom templates from ContentDir/templates to override default templates
	customTemplatesDir := filepath.Join(cfg.ContentDir, "templates")
	primeDirectory(customTemplatesDir)

	if _, err := os.Stat(customTemplatesDir); err == nil {
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
	}

	entry, _, err := content.ParseEntry(data)
	if err != nil {
		log.Println("Error parsing entry:", err)
	}

	if cfg.Mode == "serve" {
		entry.FileName = strings.Replace(filename, ".md", ".html", 1)
	} else {
		entry.FileName = filename
	}

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
		primeDirectory(filepath.Join(cfg.OutputDir, "css"))

		// render css from template and write to outputdir/css/style.css
		wr, err := os.Create(cfg.OutputDir + "/css" + "/style.css")
		if err != nil {
			log.Println("Error creating file:", err)
		}

		err = templates.ExecuteTemplate(wr, "styleCSS", cfg.Site)
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

// basic auth handler for admin pages
func basicAuthHandler(w http.ResponseWriter, r *http.Request) bool {
	if cfg.AdminUser == "" || cfg.AdminPass == "" {
		return true
	}

	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(cfg.AdminUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(cfg.AdminPass)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized\n"))
		return false
	}

	return true
}
