package contacts

import (
	"strings"
)

type Contacts []Contact

func (c *Contacts) Add(contact Contact) {
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
		if strings.Contains(contact.First, search) || strings.Contains(contact.Last, search) {
			results = append(results, contact)
		}
	}

	return results
}
