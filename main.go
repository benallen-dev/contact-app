package main

import (
	"log"
	"net/http"

	"github.com/benallen-dev/contact-app/pkg/contacts"
	"github.com/benallen-dev/contact-app/handlers"
)

const PORT = "3000"

func main() {
	var contactList contacts.Contacts

	readErr := contactList.ReadAll()
	if (readErr != nil) {
		log.Println(readErr)
	}

	http.HandleFunc("GET /", handlers.Root)

	http.HandleFunc("GET /contacts", handlers.GetContacts(&contactList))
	http.HandleFunc("GET /contacts/new", handlers.GetNewContactForm)

	http.HandleFunc("GET /contacts/{contactId}", handlers.GetContactDetails(&contactList))
	http.HandleFunc("GET /contacts/{contactId}/edit", handlers.GetEditContactForm(&contactList))

	http.HandleFunc("POST /contacts/new", handlers.PostNewContactForm(&contactList))
	http.HandleFunc("POST /contacts/{contactId}/edit", handlers.PostEditContactForm(&contactList))
	http.HandleFunc("POST /contacts/{contactId}/delete", handlers.PostDeleteContact(&contactList))

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
