package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

	// create bcrypt hash from password string
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Insert users data into "users" table
	_, err = m.DB.Exec(stmt, name, email, string(hash))
	if err != nil {
		// Use the errors.As() function to check whether the error has the type *mysql.MySQLError. 
		// If yes => error assigned to the mySQLError variable. Check if error relates to our users_uc_email key by 
		// checking if the error code equals 1062 and the contents of the error message string. 
		// If it does, we return an ErrDuplicateEmail error. 
		var mySQLError *mysql.MySQLError 
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			} 
		} 
		return err
	}
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
