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

// To Get Single Record SQL Queries || To fetch a specific snippet by ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() and id = ?`
	// To execute the SQL statement withe the QueryRow method on the connection pool
	row := m.DB.QueryRow(stmt, id)
	// To initialize a pointer to a new zeroed snippet struct
	s := &models.Snippet{}
	// To copy the values from each field in sql.Row to the corresponding field
	err := row.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// If everything goes OK then return the Snippet object
	return s, nil
}

// To return multiple record SQL Queries || to return the most recently created ten snippets, as long as they haven't expired
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// To ensure the sql.Rows result set is always properly closed before the Latest() method returns
	defer rows.Close()
	// To inialize an empty slice to hold the models.Snippets objects
	snippets := []*models.Snippet{}
	// To iterate through the roes in the result set
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets
		snippets = append(snippets, s)
	}
	// When the rows.Next() loop finished we call rows.Err() to retrieve any error during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything is OK then return the Snippets slice
	return snippets, nil
}
