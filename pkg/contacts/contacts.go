package contacts

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

var filename string = "data/contacts.csv"

// The ID for a contact is its position in the contacts array
type Contacts []Contact

func (c *Contacts) Add(contact Contact) {
	newId := len(*c)
	contact.ID = newId

	contact.Errors = make(map[string]string)

	*c = append(*c, contact)

	// TODO: Append contact to CSV
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

	// Open file on disk
	dataFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// Write to file
	csvWriter := csv.NewWriter(dataFile)
	csvWriter.WriteAll(entries)

	return nil
}

func (c *Contacts) ReadAll() error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a csv reader
	csvReader := csv.NewReader(file)
	csvData, err := csvReader.ReadAll()

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

	all := c.All()[0:id]
	newContacts := append(all[0:id], all[id:len(all)-1]...)

	c = &newContacts
}

func (c *Contacts) Update(id int, first, last, email, phone string) (Contact, error) {
	target, err := c.Get(id)
	if err != nil {
		return EmptyContact(), err
	}

	target.First = first
	target.Last = last
	target.Email = email
	target.Phone = phone

	(*c)[id] = target

	return target, nil

}

