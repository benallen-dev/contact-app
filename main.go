package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"github.com/benallen-dev/contact-app/pkg/contacts"
	"github.com/benallen-dev/contact-app/views"
)

const PORT = "3000"

func main() {
	var contactList contacts.Contacts

	contactList.ReadAll()

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


	http.HandleFunc("GET /contacts/{contactId}", func(w http.ResponseWriter, r *http.Request) {
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

		templ.Handler(views.ContactDetail(contact)).ServeHTTP(w, r)
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
