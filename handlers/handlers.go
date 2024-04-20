package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"github.com/benallen-dev/contact-app/pkg/contacts"
	"github.com/benallen-dev/contact-app/pkg/flash"
	"github.com/benallen-dev/contact-app/views"
)

var (
	contactList contacts.Contacts
)

func init() {
	err := contactList.ReadAll()
	if err != nil {
		log.Fatalf("Failed to initialize contact handlers: %v", err)
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	} else {
		http.ServeFile(w, r, "static"+r.URL.Path)
	}
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	var contacts_set []contacts.Contact

	if q == "" {
		contacts_set = contactList.All()
	} else {
		contacts_set = contactList.Search(q)
	}

	templ.Handler(views.Contacts(contacts_set, q)).ServeHTTP(w, r)
}

func GetContactDetails(w http.ResponseWriter, r *http.Request) {
	contactId, err := strconv.Atoi(r.PathValue("contactId"))
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	contact, err := contactList.Get(contactId)
	if err != nil {
		flash.Queue("Contact not found")
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}

	templ.Handler(views.ContactDetail(contact)).ServeHTTP(w, r)
}

func GetNewContactForm(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.ContactForm(contacts.NewContact("", "", "", ""))).ServeHTTP(w, r)
}

// Because contacts aren't mutated here it doesn't strictly need to be a pointer I think
func GetEditContactForm(w http.ResponseWriter, r *http.Request) {

	contactId, err := strconv.Atoi(r.PathValue("contactId"))
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	contact, err := contactList.Get(contactId)
	if err != nil {
		flash.Queue("Contact not found")
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}

	templ.Handler(views.ContactForm(contact)).ServeHTTP(w, r)
}

func PostNewContactForm(w http.ResponseWriter, r *http.Request) {
	first := r.FormValue("first_name")
	last := r.FormValue("last_name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	contact := contacts.NewContact(first, last, phone, email)
	if !contact.Validate() {
		log.Println("Invalid contact", contact.Errors)
		templ.Handler(views.ContactForm(contact)).ServeHTTP(w, r)
	}

	err := contactList.AddAndWrite(contact)
	if err != nil {
		// Todo: flash
		log.Println("Error adding contact", err)
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func PostEditContactForm(w http.ResponseWriter, r *http.Request) {
	contactId, err := strconv.Atoi(r.PathValue("contactId"))
	if err != nil {
		log.Printf("Invalid contact ID: %v", err)
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
	}

	first := r.FormValue("first_name")
	last := r.FormValue("last_name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	// Update return s err if contact is not found or fails validation
	contact, err := contactList.Update(contactId, first, last, email, phone)
	if err != nil {
		if _, ok := err.(contacts.ValidationError); ok {
			templ.Handler(views.ContactForm(contact)).ServeHTTP(w, r)
			return
		}

		flash.Queue(err.Error())
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func PostDeleteContact(w http.ResponseWriter, r *http.Request) {
	contactId, convErr := strconv.Atoi(r.PathValue("contactId"))
	if convErr != nil {
		flash.Queue("Error during delete", convErr.Error())
	}

	log.Println("Deleting contact with ID: ", contactId)
	contactList.Delete(contactId)

	writeErr := contactList.WriteAll()
	if writeErr != nil {
		flash.Queue("Error deleting contact")
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}
