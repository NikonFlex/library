package controller

import (
	"github.com/google/uuid"
	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (impl *implementation) GetAuthorBooks(request *library.GetAuthorBooksRequest, stream library.Library_GetAuthorBooksServer) error {
	impl.logger.Info("GetAuthorBooks controller: started")
	defer impl.logger.Info("GetAuthorBooks controller: finished")

	if err := request.ValidateAll(); err != nil {
		impl.logger.Error("GetAuthorBooks controller: invalid argument")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	authorID := uuid.Must(uuid.Parse(request.GetAuthorId()))
	books, err := impl.booksUseCase.GetAuthorBooks(stream.Context(), authorID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, book := range books {
		if err := stream.Send(&library.Book{
			Id:        book.ID.String(),
			Name:      book.Name,
			AuthorId:  book.AuthorIDs.Strings(),
			CreatedAt: timestamppb.New(book.CreatedAt),
			UpdatedAt: timestamppb.New(book.UpdatedAt),
		}); err != nil {
			impl.logger.Error("GetAuthorBooks controller: failed to send book - " + err.Error())
			return err
		}
	}

	return nil
}
