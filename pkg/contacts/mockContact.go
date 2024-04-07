package contacts

import (
	"log"
	"os"
	"bufio"
	"strings"
)

func CreateMockContacts() Contacts {

	file, err := os.Open("./resources/mock-contacts.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	contacts := Contacts{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		contact := Contact{
			First: parts[0],
			Last:  parts[1],
			Email: parts[2],
			Phone: parts[3],
		}

		contacts.Add(contact)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return contacts
}
