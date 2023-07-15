package main

import (
	"log"
	"net/http"
	"os"
)

func serveStaticFiles() {
	log.Println("Serving static files from", cfg.ContentDir+"/static")
	// check if static dir exists. if not, create it
	if _, err := os.Stat(cfg.ContentDir + "/static"); os.IsNotExist(err) {
		os.Mkdir(cfg.ContentDir+"/static", 0755)
	}
	// serve static files

	fs := http.FileServer(http.Dir(cfg.ContentDir + "/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func serveCSSTemplate() {
	log.Println("Serving CSS from templates")
	http.HandleFunc("/css/style.css", func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Set("Content-Type", "text/css")
		err := templates.ExecuteTemplate(wr, "styleCSS", cfg.Site.Theme)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
		}
	})
}
