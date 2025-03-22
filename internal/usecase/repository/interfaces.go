package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
)

//go:generate ../../../bin/mockgen -source=interfaces.go -destination=../../../mocks-generated/repositories_mock.go -package=mocks_generated

type BooksRepository interface {
	AddBook(ctx context.Context, book entity.Book) (entity.Book, error)
	GetBook(ctx context.Context, bookID uuid.UUID) (entity.Book, error)
	UpdateBook(ctx context.Context, bookID uuid.UUID, name string, authorIDs uuid.UUIDs) (entity.Book, error)
	GetAuthorBooks(ctx context.Context, authorID uuid.UUID) ([]entity.Book, error)
}

type AuthorsRepository interface {
	AddAuthor(ctx context.Context, author entity.Author) (entity.Author, error)
	GetAuthor(ctx context.Context, authorID uuid.UUID) (entity.Author, error)
	UpdateAuthor(ctx context.Context, authorID uuid.UUID, name string) (entity.Author, error)
}
