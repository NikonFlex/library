package controller

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *implementation) AddBook(ctx context.Context, request *library.AddBookRequest) (*library.AddBookResponse, error) {
	impl.logger.Info("AddBook controller: started")
	defer impl.logger.Info("AddBook controller: finished")

	if err := request.ValidateAll(); err != nil {
		impl.logger.Error("AddBook controller: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	book, err := impl.booksUseCase.AddBook(ctx, request.Name, impl.stringsToUUIDs(request.AuthorIds))

	impl.logger.Info("AddBook controller: TIME " + book.CreatedAt.String())
	if err != nil {
		impl.logger.Error("AddBook controller: " + err.Error())
		return nil, impl.libraryToGrpcErr(err)
	}
	impl.logger.Info("AddBook controller: TIME " + book.CreatedAt.String())

	return &library.AddBookResponse{
		Book: &library.Book{
			Id:        book.ID.String(),
			Name:      book.Name,
			AuthorId:  book.AuthorIDs.Strings(),
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
		},
	}, nil
}
