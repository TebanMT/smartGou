package domain

import "time"

type User struct {
	ID             int
	Username       string
	Name           string
	LastName       string
	SecondLastName string
	Email          string
	DailingCode    string
	Phone          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
