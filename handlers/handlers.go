package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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
	const PAGE_SIZE int = 10

	q := r.URL.Query().Get("q")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}
	
	var contacts_set []contacts.Contact

	if q == "" {
		contacts_set, err = contactList.Partial(page * PAGE_SIZE, (page + 1) * PAGE_SIZE)
		if err != nil {
			flash.Queue("Invalid page number, displaying first" + strconv.Itoa(PAGE_SIZE))
			contacts_set, _ = contactList.Partial(0, PAGE_SIZE)
		}
	} else {
		contacts_set = contactList.Search(q)
	}

	log.Println("Special token: ", r.Header.Get("X-Special-Token"))

	if r.Header.Get("HX-Trigger") == "search" {
		templ.Handler(views.ContactList(contacts_set)).ServeHTTP(w, r)
	} else {
		templ.Handler(views.Contacts(contacts_set, q, page)).ServeHTTP(w, r)
	}
}

func GetContactCount(w http.ResponseWriter, r *http.Request) {
	count := contactList.Count()

	time.Sleep(1500 * time.Millisecond)

	w.Write([]byte("(" + strconv.Itoa(count) + " total contacts)"))
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

func ValidateEmail(w http.ResponseWriter, r *http.Request) {
	contactId, err := strconv.Atoi(r.PathValue("contactId"))
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	email := r.URL.Query().Get("email")

	targetContact, err := contactList.Get(contactId)
	targetContact.Email = email
	contactList.Validate(&targetContact)

	w.Write([]byte(targetContact.Errors["email"]))
}

func PostNewContactForm(w http.ResponseWriter, r *http.Request) {
	first := r.FormValue("first_name")
	last := r.FormValue("last_name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	contact := contacts.NewContact(first, last, phone, email)
	if contactList.Validate(&contact) {
		err := contactList.AddAndWrite(contact)
		if err != nil {
			flash.Queue("Error adding contact")
			log.Println("Error adding contact", err)
		}
	} else {
		log.Println("Invalid contact", contact.Errors)
		templ.Handler(views.ContactForm(contact)).ServeHTTP(w, r)
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

	http.Redirect(w, r, "/contacts", http.StatusFound)
}


func DeleteContact(w http.ResponseWriter, r *http.Request) {
	contactId, convErr := strconv.Atoi(r.PathValue("contactId"))
	if convErr != nil {
		flash.Queue("Error during delete", convErr.Error())
	}

	contactList.Delete(contactId)

	writeErr := contactList.WriteAll()
	if writeErr != nil {
		flash.Queue("Error deleting contact")
	}

	// if from inline delete send back empty body
	if r.Header.Get("HX-Trigger") == "delete-btn" {
		w.Write([]byte(""))
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func DeleteMultipleContacts(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	parsed := parseFormBody(string(body))
	contactIds := parsed["selected_contact_ids"]

	if len(contactIds) == 0 {
		flash.Queue("No contacts selected")
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}

	log.Println("Deleting contacts with IDs: ", contactIds)
	
	for _, id := range contactIds {
		contactId, convErr := strconv.Atoi(id)
		if convErr != nil {
			flash.Queue("Error during delete: "+id, convErr.Error())
			continue
		}

		contactList.Delete(contactId)
	}

	writeErr := contactList.WriteAll()
	if writeErr != nil {
		flash.Queue("Error saving data after deleting contact")
	}

	f := "Deleted contacts: "
	for i, id := range contactIds {
		f += id
		if i < len(contactIds) - 1 {
			f += ", "
		}		
	}

	flash.Queue(f)


	contacts_set, err := contactList.Partial(0, 10)
	if err != nil {
		flash.Queue("Error fetching first 10 contacts (try reloading the page?)")
		contacts_set = []contacts.Contact{}
	}

	templ.Handler(views.Contacts(contacts_set, "", 0)).ServeHTTP(w, r)

}
