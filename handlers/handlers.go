package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"github.com/benallen-dev/contact-app/pkg/contacts"
	"github.com/benallen-dev/contact-app/views" 
)

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	} else {
		http.ServeFile(w, r, "static"+r.URL.Path)
	}
}

func GetContacts(contactList *contacts.Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		var contacts_set []contacts.Contact

		if q == "" {
			contacts_set = contactList.All()
		} else {
			contacts_set = contactList.Search(q)
		}

		templ.Handler(views.Contacts(contacts_set, q)).ServeHTTP(w, r)
	}
}

func GetContactDetails(contactList *contacts.Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func GetNewContactForm(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.ContactForm(contacts.NewContact("Ben", "Allen", "+31 6 3639 1152", "ben@benallen.dev"))).ServeHTTP(w, r)
}

func GetEditContactForm(w http.ResponseWriter, r *http.Request, contactList contacts.Contacts) {
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
}

func PostNewContactForm(w http.ResponseWriter, r *http.Request, contactList contacts.Contacts) {
	first := r.FormValue("first_name")
	last := r.FormValue("last_name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	contact := contacts.NewContact(first, last, phone, email)
	contactList.Add(contact)

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func PostEditContactForm(w http.ResponseWriter, r *http.Request, contactList contacts.Contacts) {
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
}
