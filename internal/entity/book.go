package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID        uuid.UUID
	Name      string
	AuthorIDs uuid.UUIDs
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrBookNotFound      = errors.New("Book not found")
	ErrBookAlreadyExists = errors.New("Book already exists")
)
