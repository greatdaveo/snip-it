package mysql

import (
	"database/sql"
	"snippet-box/pkg/models"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// To add a new record to the users table in the DB
func (m *UserModel) Insert(name, email, password string) error {
	// To Create a bcrypt hash of the password text
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// To check if the error is a MySQL-specific error.
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			// To check if the error code is 1062 (Duplicate entry) and if the error message contains 'u' (for unique email constraint).
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users.uc_users_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

// To verify if user exist, and return user ID if user exist
func (m *UserModel) Aunthenticate(email, password string) (int, error) {
	return 0, nil
}

// To get the details of a specific user
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
