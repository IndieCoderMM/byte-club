package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (title, content, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) List(ctx context.Context) ([]Post, error) {
	query := `SELECT id, title, content, user_id, tags, created_at, updated_at FROM posts`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}

	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostStore) Get(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, title, content, user_id, tags, version created_at, updated_at FROM posts WHERE id = $1`

	post := &Post{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.Version, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
  UPDATE posts SET title = $1, content = $2, tags = $3, version = version + 1 WHERE id = $4 AND version = $5 RETURNING version
  `

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	row, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrRecordNotFound
	}

	return nil
}
