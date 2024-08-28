package models

import (
	"database/sql"
	"time"
)

type Post struct {
	Id      int
	Title   string
	Content string
	Created time.Time
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content string) (int, error) {
	stmt := `INSERT INTO posts (title, content, created)
	VALUES(?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
