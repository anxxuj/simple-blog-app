package models

import (
	"database/sql"
	"errors"
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

func (m *PostModel) Get(id int) (*Post, error) {
	stmt := "SELECT id, title, content, created FROM posts WHERE id = ?"

	row := m.DB.QueryRow(stmt, id)

	post := &Post{}

	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return post, nil
}
