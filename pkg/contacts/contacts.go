package contacts

import (
	"strings"
	"errors"
)

type Contacts []Contact

// Not a pointer so arguments are passed by value
// This means we can mutate them freely
func (c *Contacts) Add(contact Contact) {
	newId := len(*c)
	contact.ID = newId

	contact.Errors = make(map[string]string)

	*c = append(*c, contact)
}

func (c *Contacts) All() Contacts {
	return *c
}

func (c *Contacts) Search(search string) Contacts {
	if search == "" {
		return *c
	}

	var results Contacts
	for _, contact := range *c {
		match := strings.ToLower(search)

		firstMatch := strings.Contains(strings.ToLower(contact.First), match)
		lastMatch := strings.Contains(strings.ToLower(contact.Last), match)
		emailMatch := strings.Contains(strings.ToLower(contact.Email), match)

		if firstMatch || lastMatch || emailMatch {
			results = append(results, contact)
		}
	}

	return results
}

func (c *Contacts) Get(id int) (Contact, error) {
	if id < 0 || id >= len(*c) {
		return Contact{}, errors.New("Contact not found")
	}

	return (*c)[id], nil
}
