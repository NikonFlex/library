package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *implementation) UpdateBook(ctx context.Context, request *library.UpdateBookRequest) (*library.UpdateBookResponse, error) {
	impl.logger.Info("UpdateBook controller: started")
	defer impl.logger.Info("UpdateBook controller: finished")

	if err := request.ValidateAll(); err != nil {
		impl.logger.Error("UpdateBook controller: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	bookID := uuid.Must(uuid.Parse(request.GetId()))
	_, err := impl.booksUseCase.UpdateBook(ctx, bookID, request.GetName(), impl.stringsToUUIDs(request.GetAuthorIds()))
	if err != nil {
		impl.logger.Error("UpdateBook controller: " + err.Error())
		return nil, impl.libraryToGrpcErr(err)
	}

	return &library.UpdateBookResponse{}, nil
}
