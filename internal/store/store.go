package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound    = errors.New("Record not found")
	QueryTimeoutDuration = time.Second * 5
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		List(context.Context) ([]Post, error)
		Get(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
