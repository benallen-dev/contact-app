package main

import (
	"log"
	"strconv"
	"net/http"

	"github.com/a-h/templ"

	"github.com/benallen-dev/contact-app/pkg/contacts"
	"github.com/benallen-dev/contact-app/views"
)

const PORT = "3000"

func main() {
	contactList := contacts.CreateMockContacts()

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		} else {
			http.ServeFile(w, r, "static"+r.URL.Path)
		}
	})

	http.HandleFunc("GET /contacts", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		var contacts_set []contacts.Contact

		if q == "" {
			contacts_set = contactList.All()
		} else {
			contacts_set = contactList.Search(q)
		}

		templ.Handler(views.Contacts(contacts_set, q)).ServeHTTP(w, r)
	})

	http.HandleFunc("GET /contacts/new", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.ContactForm(contacts.NewContact("", "", "", ""))).ServeHTTP(w, r)
	})

	http.HandleFunc("POST /contacts/new", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented"))
	})


	http.HandleFunc("GET /contacts/{contactId}/edit", func(w http.ResponseWriter, r *http.Request) {
		contactId, err := strconv.Atoi(r.PathValue("contactId"))
		if err != nil {
			http.Error(w, "Invalid contact ID", http.StatusBadRequest)
			return
		}

		contact, err := contactList.Get(contactId)
		if err != nil {
			http.Error(w, "Contact not found", http.StatusNotFound)
			return
		}

		templ.Handler(views.ContactForm(contact)).ServeHTTP(w, r)
	})


	http.HandleFunc("POST /contacts/{contactId}/edit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented"))
	})

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
