package api

import "net/http"

// InitRoutes initializes the API routes.
func InitRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
	http.HandleFunc("/start", StartProcessing)
	http.HandleFunc("/stats", Stats)
}
