package library

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
)

func (library *libraryImpl) AddAuthor(ctx context.Context, name string) (entity.Author, error) {
	library.logger.Info("AddAuthor use case: started")
	defer library.logger.Info("AddAuthor use case: finished")
	return library.authorsRepository.AddAuthor(ctx, entity.Author{
		Name: name,
	})
}

func (library *libraryImpl) GetAuthor(ctx context.Context, authorID uuid.UUID) (entity.Author, error) {
	library.logger.Info("GetAuthor use case: started")
	defer library.logger.Info("GetAuthor use case: finished")
	return library.authorsRepository.GetAuthor(ctx, authorID)
}

func (library *libraryImpl) UpdateAuthor(ctx context.Context, authorID uuid.UUID, name string) (entity.Author, error) {
	library.logger.Info("UpdateAuthor use case: started")
	defer library.logger.Info("GetAuthor use case: finished")
	return library.authorsRepository.UpdateAuthor(ctx, authorID, name)
}
