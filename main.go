package main


import (
	"net/http"
	"log"
)

const PORT = "3000"

func main() {

	// Just using default mux for simplicity
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		} else {
			http.ServeFile(w, r, "static"+r.URL.Path)
		}
	})
	
	http.HandleFunc("GET /contacts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Contacts page"))
	})
	
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
