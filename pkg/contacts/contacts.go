package contacts

import (
	"errors"
	"strings"
)

// The ID for a contact is its position in the contacts array
type Contacts []Contact

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

// This is kinda meh because it rewrites the entire file
func (c *Contacts) WriteAll() error {
	entries := [][]string{}

	// Convert the list of contacts to a list of comma-seperated values
	for _, contact := range c.All() {
		entries = append(entries, contact.ToStringArray())
	}

	err := writeCsv(filename, entries)
	if err != nil {
		return err
	}

	return nil
}

func (c *Contacts) ReadAll() error {
	csvData, err := readCsv(filename)
	if err != nil {
		return err
	}

	for _, entry := range csvData {
		newContact, err := NewContactFromCSV(entry...)
		if err != nil {
			return err
		}

		c.Add(newContact)
	}

	return nil
}

func (c *Contacts) Delete(id int) {
	// I should have used a map, but I was not smart and so
	// here we are slicing arrays for each delete

	all := c.All()

	newContacts := append(all[0:id], all[id:len(all)-1]...)

	*c = newContacts

// 	err := (*c).WriteAll()
// 	if err != nil {
// 		return err
// 	}

	return
}

func (c *Contacts) Update(id int, first, last, email, phone string) (Contact, error) {
	if id < 0 || id >= len(*c) {
		return EmptyContact(), errors.New("Contact not found")
	}

	target := &(*c)[id]

	target.First = first
	target.Last = last
	target.Email = email
	target.Phone = phone

	(*c)[id] = *target

	err := c.WriteAll()
	if err != nil {
		return EmptyContact(), err
	}

	return *target, nil

}
