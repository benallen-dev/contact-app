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
