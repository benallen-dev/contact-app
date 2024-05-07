package main

import (
	"log"
	"net/http"

	"github.com/benallen-dev/contact-app/handlers"
)

const PORT = "3000"

func main() {
	http.HandleFunc("GET /", handlers.Root)

	http.HandleFunc("GET /contacts", handlers.GetContacts)
	http.HandleFunc("GET /contacts/count", handlers.GetContactCount)
	http.HandleFunc("GET /contacts/new", handlers.GetNewContactForm)

	http.HandleFunc("GET /contacts/{contactId}", handlers.GetContactDetails)
	http.HandleFunc("GET /contacts/{contactId}/edit", handlers.GetEditContactForm)
	http.HandleFunc("GET /contacts/{contactId}/email", handlers.ValidateEmail)

	http.HandleFunc("POST /contacts/new", handlers.PostNewContactForm)
	http.HandleFunc("POST /contacts/{contactId}/edit", handlers.PostEditContactForm)
	http.HandleFunc("POST /contacts/{contactId}/delete", handlers.PostDeleteContact)

	// Fancy HTMX endpoints here
	http.HandleFunc("DELETE /contacts/{contactId}", handlers.DeleteContact)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
