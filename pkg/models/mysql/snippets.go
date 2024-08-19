package mysql

import (
	"database/sql"
	"snippet-box/pkg/models"
)

// To define a SnippetModel type that wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// To insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// To return a specific snippet based on its id
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}

