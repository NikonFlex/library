package library

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
)

func (library *libraryImpl) AddBook(ctx context.Context, name string, authorIDs uuid.UUIDs) (entity.Book, error) {
	library.logger.Info("AddBook use case: started")
	defer library.logger.Info("AddBook use case: finished")

	return library.booksRepository.AddBook(ctx, entity.Book{
		Name:      name,
		AuthorIDs: authorIDs,
	})
}

func (library *libraryImpl) GetBook(ctx context.Context, bookID uuid.UUID) (entity.Book, error) {
	library.logger.Info("GetBook use case: started")
	defer library.logger.Info("GetBook use case: finished")

	return library.booksRepository.GetBook(ctx, bookID)
}

func (library *libraryImpl) UpdateBook(ctx context.Context, bookID uuid.UUID, name string, authorIDs uuid.UUIDs) (entity.Book, error) {
	library.logger.Info("UpdateBook use case: started")
	defer library.logger.Info("UpdateBook use case: finished")

	return library.booksRepository.UpdateBook(ctx, bookID, name, authorIDs)
}

func (library *libraryImpl) GetAuthorBooks(ctx context.Context, authorID uuid.UUID) ([]entity.Book, error) {
	library.logger.Info("GetAuthorBooks use case: started")
	defer library.logger.Info("GetAuthorBooks use case: finished")

	return library.booksRepository.GetAuthorBooks(ctx, authorID)
}
