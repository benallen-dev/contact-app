package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"

	"github.com/benallen-dev/contact-app/pkg/contacts"
	"github.com/benallen-dev/contact-app/views"
)

const PORT = "3000"

func newFilledContacts() contacts.Contacts {
	return contacts.Contacts{
		contacts.Contact{
			ID:    1,
			First: "Ben",
			Last:  "Allen",
			Phone: "123-456-7890",
			Email: "foo@bar.baz",
		},
		contacts.Contact{
			ID:    2,
			First: "John",
			Last:  "Doe",
			Phone: "098-765-4321",
			Email: "johndoe@random.biz",
			},
	}
}
	

func main() {
	contactList := newFilledContacts()

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

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
