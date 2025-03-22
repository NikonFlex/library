package controller

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/project/library/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (impl *implementation) libraryToGrpcErr(err error) error {
	switch {
	case errors.Is(err, entity.ErrBookNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entity.ErrBookAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entity.ErrAuthorNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entity.ErrAuthorAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// call ONLY after ValidateAll
func (impl *implementation) stringsToUUIDs(strings []string) uuid.UUIDs {
	uuids := make([]uuid.UUID, 0, len(strings))
	for _, s := range strings {
		u := uuid.Must(uuid.Parse(s))
		uuids = append(uuids, u)
	}
	return uuids
}
