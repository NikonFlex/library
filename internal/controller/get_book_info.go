package controller

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *implementation) GetBookInfo(ctx context.Context, request *library.GetBookInfoRequest) (*library.GetBookInfoResponse, error) {
	impl.logger.Info("GetBookInfo controller: started")
	defer impl.logger.Info("GetBookInfo controller: finished")

	if err := request.ValidateAll(); err != nil {
		impl.logger.Error("GetBookInfo controller: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	bookID := uuid.Must(uuid.Parse(request.GetId()))
	book, err := impl.booksUseCase.GetBook(ctx, bookID)
	if err != nil {
		impl.logger.Error("GetBookInfo controller: " + err.Error())
		return nil, impl.libraryToGrpcErr(err)
	}

	impl.logger.Info("AddBook controller: TIME " + book.CreatedAt.String())
	return &library.GetBookInfoResponse{
		Book: &library.Book{
			Id:        book.ID.String(),
			Name:      book.Name,
			AuthorId:  book.AuthorIDs.Strings(),
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
		},
	}, nil
}
