package controller

import (
	"context"

	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *implementation) RegisterAuthor(ctx context.Context, request *library.RegisterAuthorRequest) (*library.RegisterAuthorResponse, error) {
	impl.logger.Info("RegisterAuthor controller: started")
	defer impl.logger.Info("RegisterAuthor controller: finished")

	if err := request.ValidateAll(); err != nil {
		impl.logger.Error("RegisterAuthor controller: invalid argument")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	author, err := impl.authorsUseCase.AddAuthor(ctx, request.Name)
	if err != nil {
		impl.logger.Error("RegisterAuthor controller: " + err.Error())
		return nil, impl.libraryToGrpcErr(err)
	}

	return &library.RegisterAuthorResponse{
		Id: author.ID.String(),
	}, nil
}
