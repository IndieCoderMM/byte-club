package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("Record not found")
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		Get(context.Context, int64) (*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}
