package library

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
	"github.com/project/library/internal/usecase/repository"
	"go.uber.org/zap"
)

//go:generate ../../../bin/mockgen -source=interfaces.go -destination=../../../mocks-generated/usecases_mock.go -package=mocks_generated

type BooksUseCase interface {
	AddBook(ctx context.Context, name string, authorIDs uuid.UUIDs) (entity.Book, error)
	GetBook(ctx context.Context, bookID uuid.UUID) (entity.Book, error)
	UpdateBook(ctx context.Context, bookID uuid.UUID, name string, authorIDs uuid.UUIDs) (entity.Book, error)
	GetAuthorBooks(ctx context.Context, authorID uuid.UUID) ([]entity.Book, error)
}

type AuthorsUseCase interface {
	AddAuthor(ctx context.Context, name string) (entity.Author, error)
	GetAuthor(ctx context.Context, authorID uuid.UUID) (entity.Author, error)
	UpdateAuthor(ctx context.Context, authorID uuid.UUID, name string) (entity.Author, error)
}

var _ BooksUseCase = (*libraryImpl)(nil)
var _ AuthorsUseCase = (*libraryImpl)(nil)

type libraryImpl struct {
	logger            *zap.Logger
	booksRepository   repository.BooksRepository
	authorsRepository repository.AuthorsRepository
}

func New(
	logger *zap.Logger,
	booksRepository repository.BooksRepository,
	authorsRepository repository.AuthorsRepository,
) *libraryImpl {
	return &libraryImpl{
		logger:            logger,
		booksRepository:   booksRepository,
		authorsRepository: authorsRepository,
	}
}
