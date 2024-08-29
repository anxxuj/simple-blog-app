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

func (m *PostModel) GetAll() ([]*Post, error) {
	stmt := "SELECT id, title, content, created FROM posts ORDER BY id DESC"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}

	for rows.Next() {
		post := &Post{}

		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Created)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (m *PostModel) Update(postId int, title, content string) error {
	stmt := `UPDATE posts
	SET title = ?, content = ?, created = UTC_TIMESTAMP()
	WHERE id = ?`

	_, err := m.DB.Exec(stmt, title, content, postId)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) Delete(id int) error {
	stmt := "DELETE FROM posts WHERE id = ?"

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
