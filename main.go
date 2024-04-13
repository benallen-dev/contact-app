package main

import (
	"log"
	"net/http"
	"strconv"

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

	http.HandleFunc("GET /contacts/{contactId}/edit", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetEditContactForm(w, r, contactList)
	})

	http.HandleFunc("POST /contacts/new", func(w http.ResponseWriter, r *http.Request) {
		first := r.FormValue("first_name")
		last := r.FormValue("last_name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		contact := contacts.NewContact(first, last, phone, email)
		contactList.Add(contact)

		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})

	http.HandleFunc("POST /contacts/{contactId}/edit", func(w http.ResponseWriter, r *http.Request) {
		contactId, err := strconv.Atoi(r.PathValue("contactId"))
		if err != nil {
			log.Printf("Invalid contact ID: %v", err)
			http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		}

		first := r.FormValue("first_name")
		last := r.FormValue("last_name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		contactList.Update(contactId, first, last, email, phone)

		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
