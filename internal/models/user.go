package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Wraps connection pool ?
type UserModel struct {
	DB *sql.DB
}

// Add user record
func (m *UserModel) Insert(name, email, password string) error { 
	return nil 
}

// Verify if a user with the provided "email" & "password" exists
// Return user ID on success
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 1, nil
}

// Check if user with ID exists
// Return bool
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}