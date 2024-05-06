package contacts

import (
	"errors"
	"log"
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

func (c *Contacts) AddAndWrite(contact Contact) error {
	c.Add(contact)

	err := c.WriteAll()
	if err != nil {
		return err
	}

	return nil
}

func (c *Contacts) All() Contacts {
	return *c
}

func (c *Contacts) Count() int {
	return len(*c)
}

func (c *Contacts) Partial (start, end int) (Contacts, error) {
	if start < 0 || end < 0 || start > end {
		return Contacts{}, errors.New("Invalid range")
	}

	if end > len(*c) {
		end = len(*c)
	}

	return (*c)[start:end], nil
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

func (c *Contacts) Delete(id int) error {
	// I should have used a map, but I was not smart and so
	// here we are slicing arrays for each delete

	all := c.All()

	newContacts := append(all[0:id], all[id:len(all)-1]...)

	*c = newContacts

	err := (*c).WriteAll()
	if err != nil {
		return err
	}

	return nil
}

// Update a contact by ID
//
// If the contact is not found or otherwise fails validation, an error is
// returned. In the case of validation failure the contact is returned with
// errors filled in.
//
// If the contact is updated successfully, the updated contact is returned
func (c *Contacts) Update(id int, first, last, email, phone string) (Contact, error) {
	if id < 0 || id >= len(*c) {
		return EmptyContact(), errors.New("Contact not found")
	}

	updatedContact := NewContact(first, last, phone, email)
	updatedContact.ID = id

	if !c.Validate(&updatedContact) {
		return updatedContact, ValidationError{"Invalid contact"}
	}

	log.Printf("%+v", updatedContact)

	(*c)[id] = updatedContact

	err := c.WriteAll()
	if err != nil {
		return EmptyContact(), err
	}

	//return *target, nil
	return updatedContact, nil
}

// Validates the passed contact. Mutates contact.Errors and returns a bool indicating if there were any errors
func (c *Contacts) Validate(contact *Contact) bool {
	validContact := contact.Validate()
	uniqueEmail := c.checkEmailIsFree(contact)

	log.Print(contact.Errors)

	return validContact && uniqueEmail
}

func (c *Contacts) checkEmailIsFree(contact *Contact) bool {
	// Check if email is already in use
	for _, existing := range c.All() {
		log.Printf("Checking %s (%d) - %s (%d)", existing.Email, existing.ID, contact.Email, contact.ID)
		if existing.Email == contact.Email && existing.ID != contact.ID {

			log.Printf("%s already exists", existing.Email)
			contact.Errors["email"] += "Email already in use"
			return false
		}
	}

	return true
}
