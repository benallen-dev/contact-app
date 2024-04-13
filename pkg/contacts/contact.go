package contacts

type Contact struct {
	ID     int
	First  string
	Last   string
	Phone  string
	Email  string
	Errors map[string]string
}

func NewContact(first, last, phone, email string) Contact {
	return Contact{
		ID:		-1,
		First:  first,
		Last:   last,
		Phone:  phone,
		Email:  email,
		Errors: make(map[string]string),
	}
}

func EmptyContact() Contact {
	return Contact{
		ID:     -1,
		First:  "",
		Last:   "",
		Email:  "",
		Phone:  "",
		Errors: make(map[string]string),
	}
}

func NewContactFromCSV(d ...string) (Contact, error) {
	return Contact{
		First: d[0],
		Last:  d[1],
		Email: d[2],
		Phone: d[3],
	}, nil
}

func (c *Contact) ToStringArray() []string {
	return []string{
		c.First,
		c.Last,
		c.Email,
		c.Phone,
	}
}
