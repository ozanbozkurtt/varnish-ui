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

func endpointStatsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /endpoint_stats endpoint")

	stats, err := varnish.GetVarnishEndpointStats()
	if err != nil {
		log.Printf("Error getting endpoint stats: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(stats)
}

func main() {
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/charts", chartsHandler)
	http.HandleFunc("/endpoint_stats", endpointStatsHandler)

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
