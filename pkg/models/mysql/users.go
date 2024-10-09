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
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// To retrieve the user id and hashed password
	var id int
	var hashedPassword []byte
	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// To check if the hashed password and password text matches
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// Otherwise the password is correct, return the user id
	return id, nil
}

// To get the details of a specific user
func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
