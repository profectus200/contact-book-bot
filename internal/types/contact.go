package types

import (
	"fmt"
	"time"
)

type Contact struct {
	ContactID   int
	Name        string
	Email       string
	Phone       string
	Birthday    time.Time
	Description string
}

func NewContact() *Contact {
	return &Contact{
		ContactID:   0,
		Name:        "New contact",
		Email:       "",
		Phone:       "",
		Birthday:    time.Now(),
		Description: "",
	}
}

func (c *Contact) ToString() string {
	str := fmt.Sprintf("ID: %d\n", c.ContactID)
	str += fmt.Sprintf("Name: %s\n", c.Name)
	if c.Email != "" {
		str += fmt.Sprintf("Email: %s\n", c.Email)
	}
	if c.Phone != "" {
		str += fmt.Sprintf("Phone: %s\n", c.Phone)
	}
	if c.Birthday.Year() != time.Now().Year() {
		str += fmt.Sprintf("Birthday: %s\n", c.Birthday.Format("02.01"))
	}
	if c.Description != "" {
		str += fmt.Sprintf("Description: %s\n", c.Description)
	}
	return str
}
