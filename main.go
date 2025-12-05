package main

import (
	"log"
	"mime"
	"net/http"
	"os"
)

func main() {
	// API routes
	http.HandleFunc("/api/rooms", enableCORS(handleRooms))
	http.HandleFunc("/api/rooms/", enableCORS(handleRooms))

	// Ensure .webp files are served with the correct MIME type
	// Some Go stdlib versions don't register .webp by default.
	_ = mime.AddExtensionType(".webp", "image/webp")

	// Serve static frontend files if the static directory exists
	if _, err := os.Stat("static"); err == nil {
		fs := http.FileServer(http.Dir("static"))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Serve index.html for all non-API routes (SPA support)
			if r.URL.Path != "/" && !fileExists("static"+r.URL.Path) {
				http.ServeFile(w, r, "static/index.html")
				return
			}
			fs.ServeHTTP(w, r)
		})
		log.Println("Serving static files from ./static")
	}

	log.Println("CrossClues server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
