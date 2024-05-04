package contacts

import (
	"strings"
	"regexp"
)

type Contact struct {
	ID     int
	First  string
	Last   string
	Phone  string
	Email  string
	Errors map[string]string
}

type ValidationError struct {
	s string
}

func (e ValidationError) Error() string {
	return e.s
}

func NewContact(first, last, phone, email string) Contact {
	return Contact{
		ID:     -1,
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

func (c *Contact) Validate() bool {
	emailRegex := regexp.MustCompile("^[\\w+-\\.]+@([\\w-]+\\.)+[\\w]{2,4}$")
	var errors = map[string][]string{
		"first": {},
		"last":  {},
		"email": {},
		"phone": {},
	}

	// Clear out any previous errors
	c.Errors = make(map[string]string)

	if c.First == "" {
		errors["first"] = append(errors["first"], "First name is required")
	}

	if c.Last == "" {
		errors["last"] = append(errors["last"], "Last name is required")
	}

	if c.Email == "" {
		errors["email"] = append(errors["email"], "Email is required")
	}

	if !emailRegex.MatchString(c.Email) {
		errors["email"] = append(errors["email"], "Email is invalid - must contain @ and a tld no longer than 4 characters")
	}

	if c.Phone == "" {
		errors["phone"] = append(errors["phone"], "Phone is required")
	}

	for key,msgs := range errors {
		if len(msgs) > 0 {
			c.Errors[key] = strings.Join(msgs, ", ")
		}
	}

	return len(c.Errors) == 0
}
