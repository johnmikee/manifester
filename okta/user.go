package okta

import (
	"time"
)

// UserResponse holds information on the user returned when querying the /users endpoint
type User struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	Created         time.Time `json:"created"`
	Activated       time.Time `json:"activated"`
	StatusChanged   time.Time `json:"statusChanged"`
	LastLogin       time.Time `json:"lastLogin"`
	LastUpdated     time.Time `json:"lastUpdated"`
	PasswordChanged time.Time `json:"passwordChanged"`
	Profile         Profile   `json:"profile"`
}

// Profiles holds information on the users profile
type Profile struct {
	LastName    string `json:"lastName"`
	Manager     string `json:"manager"`
	SecondEmail string `json:"secondEmail"`
	ManagerID   string `json:"managerId"`
	Title       string `json:"title"`
	Login       string `json:"login"`
	FirstName   string `json:"firstName"`
	UserType    string `json:"userType"`
	Department  string `json:"department"`
	StartDate   string `json:"startDate"`
	Email       string `json:"email"`
}
