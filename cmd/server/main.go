package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"my-varnish-stats/internal/varnish"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("web/static/templates/stats.html"))
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /stats endpoint")

	stats, err := varnish.GetVarnishStats()
	if err != nil {
		log.Printf("Error getting varnish stats: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statsMap := make(map[string]int)
	err = json.Unmarshal(stats, &statsMap)
	if err != nil {
		log.Printf("Error parsing varnish stats: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "stats.html", statsMap)
	if err != nil {
		log.Printf("Error rendering template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func chartsHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/stats", http.StatusMovedPermanently)
}

func main() {
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/charts", chartsHandler) // Redirect /charts to /stats

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
