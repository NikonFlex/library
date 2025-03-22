package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *implementation) GetAuthorInfo(ctx context.Context, request *library.GetAuthorInfoRequest) (*library.GetAuthorInfoResponse, error) {
	impl.logger.Info("GetAuthorInfo controller: started")
	defer impl.logger.Info("GetAuthorInfo controller: finished")

	if err := request.ValidateAll(); err != nil {
		impl.logger.Error("GetAuthorInfo controller: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	authorID := uuid.Must(uuid.Parse(request.Id))
	author, err := impl.authorsUseCase.GetAuthor(ctx, authorID)
	if err != nil {
		impl.logger.Error("GetAuthorInfo controller: " + err.Error())
		return nil, impl.libraryToGrpcErr(err)
	}

	return &library.GetAuthorInfoResponse{
		Id:   author.ID.String(),
		Name: author.Name,
	}, nil
}
