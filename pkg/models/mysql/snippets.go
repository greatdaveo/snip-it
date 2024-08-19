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
// The SQL statement to be executed
stmt := `INSERT INTO snippets (title, content, created, expires) 
		 VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
// To execute the statement
result, err := m.DB.Exec(stmt, title, content, expires)
if err != nil {
	return 0, err
}
// To get the ID of the newly inserted record in the snippets table
id, err := result.LastInsertId()
if err != nil {
	return 0, err
}

	// The ID returned has the type int64, so it is converted to an int type before returning it
	return int(id), nil
}

// To return a specific snippet based on its id
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}

