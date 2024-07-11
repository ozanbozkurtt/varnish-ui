package main

import (
	"log"
	"my-varnish-stats/internal/varnish"
	"net/http"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /stats endpoint")

	stats, err := varnish.GetVarnishStats()
	if err != nil {
		log.Printf("Error getting varnish stats: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(stats)
}

func chartsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/charts.html")
}

func main() {
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/charts", chartsHandler)

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
