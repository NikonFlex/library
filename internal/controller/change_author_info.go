package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *implementation) ChangeAuthorInfo(ctx context.Context, request *library.ChangeAuthorInfoRequest) (*library.ChangeAuthorInfoResponse, error) {
	impl.logger.Info("ChangeAuthorInfo controller: started")
	defer impl.logger.Info("ChangeAuthorInfo controller: finished")

	if err := request.ValidateAll(); err != nil {
		impl.logger.Error("ChangeAuthorInfo controller: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	authorID := uuid.Must(uuid.Parse(request.GetId()))
	_, err := impl.authorsUseCase.UpdateAuthor(ctx, authorID, request.GetName())
	if err != nil {
		impl.logger.Error("ChangeAuthorInfo controller: " + err.Error())
		return nil, impl.libraryToGrpcErr(err)
	}

	return &library.ChangeAuthorInfoResponse{}, nil
}
