package entity

import (
	"errors"

	"github.com/google/uuid"
)

type Author struct {
	ID   uuid.UUID
	Name string
}

var (
	ErrAuthorNotFound      = errors.New("Author not found")
	ErrAuthorAlreadyExists = errors.New("Author already exists")
)
