package api

import "net/http"

// InitRoutes initialises the API routes
func InitRoutes() {
	status = ProcessingStatusNotStarted

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	http.HandleFunc("/start", StartHandler)
	http.HandleFunc("/stats", StatsHandler)
	http.HandleFunc("/pause", PauseHandler)
}
