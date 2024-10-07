package mysql

import (
	"database/sql"
	"snippet-box/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// To add a new record to the users table in the DB
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// To verify if user exist, and return user ID if user exist
func (m *UserModel) Aunthenticate(email, password string) (int, error) {
	return 0, nil
}

// To get the details of a specific user
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

