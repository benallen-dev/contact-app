package main

type Contact struct {
	ID int
	First string
	Last string
	Phone string
	Email string
}

type Contacts []Contact

func (c *Contacts) Add(contact Contact) {
	*c = append(*c, contact)
}

func (c *Contacts) All() Contacts {
	return *c
}
