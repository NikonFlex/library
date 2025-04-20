package controller

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/project/library/generated/api/library"
	"github.com/project/library/internal/entity"
	custommocks "github.com/project/library/mocks-custom"
	generatedmocks "github.com/project/library/mocks-generated"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestController_AddBook(t *testing.T) {
	t.Parallel()

	basicBook := entity.Book{
		ID:        uuid.Nil,
		Name:      "book",
		AuthorIDs: []uuid.UUID{uuid.Nil},
	}

	nilBook := entity.Book{}

	basicBookRequest := library.AddBookRequest{
		Name:      basicBook.Name,
		AuthorIds: basicBook.AuthorIDs.Strings(),
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().AddBook(ctx, basicBook.Name, basicBook.AuthorIDs).
			Return(basicBook, nil)

		result, err := controllerSetup.service.AddBook(ctx, &basicBookRequest)
		require.Equal(t, result.GetBook().GetId(), basicBook.ID.String())
		require.Equal(t, result.GetBook().GetName(), basicBook.Name)
		require.NoError(t, err)
	})

	t.Run("book already exist", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().AddBook(ctx, basicBook.Name, basicBook.AuthorIDs).
			Return(nilBook, entity.ErrBookAlreadyExists)

		result, err := controllerSetup.service.AddBook(ctx, &basicBookRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.AlreadyExists, entity.ErrBookAlreadyExists.Error()))
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().AddBook(ctx, basicBook.Name, basicBook.AuthorIDs).
			Return(nilBook, entity.ErrAuthorNotFound)

		result, err := controllerSetup.service.AddBook(ctx, &basicBookRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.NotFound, entity.ErrAuthorNotFound.Error()))
	})

	t.Run("author invalid uuid", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.AddBook(ctx, &library.AddBookRequest{
			Name:      basicBook.Name,
			AuthorIds: []string{"invalid-uuid"},
		})
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestController_GetBook(t *testing.T) {
	t.Parallel()

	basicBook := entity.Book{
		ID:        uuid.Nil,
		Name:      "book",
		AuthorIDs: []uuid.UUID{uuid.Nil},
	}

	nilBook := entity.Book{}

	basicBookRequest := library.GetBookInfoRequest{
		Id: basicBook.ID.String(),
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().GetBook(ctx, basicBook.ID).
			Return(basicBook, nil)

		result, err := controllerSetup.service.GetBookInfo(ctx, &basicBookRequest)
		require.Equal(t, result.GetBook().GetId(), basicBook.ID.String())
		require.Equal(t, result.GetBook().GetName(), basicBook.Name)
		require.NoError(t, err)
	})

	t.Run("book not found", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().GetBook(ctx, basicBook.ID).
			Return(nilBook, entity.ErrBookNotFound)

		result, err := controllerSetup.service.GetBookInfo(ctx, &basicBookRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.NotFound, entity.ErrBookNotFound.Error()))
	})

	t.Run("book invalid uuid", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.GetBookInfo(ctx, &library.GetBookInfoRequest{
			Id: "not uuid",
		})
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestController_UpdateBook(t *testing.T) {
	t.Parallel()

	basicBook := entity.Book{
		ID:        uuid.Nil,
		Name:      "book",
		AuthorIDs: []uuid.UUID{uuid.Nil},
	}

	nilBook := entity.Book{}

	updatedBasicBook := entity.Book{
		ID:        basicBook.ID,
		Name:      basicBook.Name + " update",
		AuthorIDs: []uuid.UUID{uuid.Nil, uuid.Max},
	}

	basicBookRequest := library.UpdateBookRequest{
		Id:        updatedBasicBook.ID.String(),
		Name:      updatedBasicBook.Name,
		AuthorIds: updatedBasicBook.AuthorIDs.Strings(),
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().UpdateBook(ctx, updatedBasicBook.ID, updatedBasicBook.Name, updatedBasicBook.AuthorIDs).
			Return(updatedBasicBook, nil)

		result, err := controllerSetup.service.UpdateBook(ctx, &basicBookRequest)
		require.NotNil(t, result)
		require.NoError(t, err)
	})

	t.Run("book not found", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().UpdateBook(ctx, updatedBasicBook.ID, updatedBasicBook.Name, updatedBasicBook.AuthorIDs).
			Return(nilBook, entity.ErrBookNotFound)

		result, err := controllerSetup.service.UpdateBook(ctx, &basicBookRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.NotFound, entity.ErrBookNotFound.Error()))
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.booksUseCaseMock.EXPECT().UpdateBook(ctx, updatedBasicBook.ID, updatedBasicBook.Name, updatedBasicBook.AuthorIDs).
			Return(nilBook, entity.ErrAuthorNotFound)

		result, err := controllerSetup.service.UpdateBook(ctx, &basicBookRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.NotFound, entity.ErrAuthorNotFound.Error()))
	})

	t.Run("book invalid uuid", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.UpdateBook(ctx, &library.UpdateBookRequest{
			Id:        "not uuid",
			Name:      basicBook.Name,
			AuthorIds: []string{},
		})
		require.Nil(t, result)
		require.Error(t, err)
	})

	t.Run("author invalid uuid", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.UpdateBook(ctx, &library.UpdateBookRequest{
			Id:        basicBook.ID.String(),
			Name:      basicBook.Name,
			AuthorIds: []string{"not uuid"},
		})
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestController_GetAuthorBooks(t *testing.T) {
	vasya := uuid.Nil

	book1 := entity.Book{ID: uuid.Nil, Name: "Book 1", AuthorIDs: []uuid.UUID{vasya}}
	book2 := entity.Book{ID: uuid.Max, Name: "Book 2", AuthorIDs: []uuid.UUID{vasya}}
	basicRequest := library.GetAuthorBooksRequest{
		AuthorId: vasya.String(),
	}

	t.Run("success", func(t *testing.T) {
		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})

		controllerSetup := createControllerSetup(controller)

		controllerSetup.booksUseCaseMock.EXPECT().GetAuthorBooks(gomock.Any(), vasya).
			Return([]entity.Book{book1, book2}, nil)

		streamMock := custommocks.NewMockLibraryGetAuthorBooksServer(nil, nil)

		err := controllerSetup.service.GetAuthorBooks(&basicRequest, streamMock)
		require.NoError(t, err)
		require.Equal(t, len(streamMock.SentBooks), 2)
	})

	t.Run("stream send error", func(t *testing.T) {
		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})

		controllerSetup := createControllerSetup(controller)

		controllerSetup.booksUseCaseMock.EXPECT().GetAuthorBooks(gomock.Any(), vasya).
			Return([]entity.Book{book1, book2}, nil)

		expectedErr := errors.New("send error")
		streamMock := custommocks.NewMockLibraryGetAuthorBooksServer(nil, expectedErr)

		err := controllerSetup.service.GetAuthorBooks(&basicRequest, streamMock)
		require.ErrorIs(t, expectedErr, err)
		require.Equal(t, len(streamMock.SentBooks), 0)
	})

	t.Run("author invalid uuid", func(t *testing.T) {
		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})

		controllerSetup := createControllerSetup(controller)
		streamMock := custommocks.NewMockLibraryGetAuthorBooksServer(nil, nil)
		err := controllerSetup.service.GetAuthorBooks(&library.GetAuthorBooksRequest{
			AuthorId: "not-uuid",
		}, streamMock)
		require.Equal(t, len(streamMock.SentBooks), 0)
		require.Error(t, err)
	})
}

func TestController_AddAuthor(t *testing.T) {
	t.Parallel()

	basicAuthor := entity.Author{
		ID:   uuid.Nil,
		Name: "author",
	}

	nilAuthor := entity.Author{}

	basicAuthorRequest := library.RegisterAuthorRequest{
		Name: basicAuthor.Name,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.authorsUseCaseMock.EXPECT().AddAuthor(ctx, basicAuthor.Name).
			Return(basicAuthor, nil)

		result, err := controllerSetup.service.RegisterAuthor(ctx, &basicAuthorRequest)
		require.Equal(t, result, &library.RegisterAuthorResponse{
			Id: basicAuthor.ID.String(),
		})
		require.NoError(t, err)
	})

	t.Run("author already exist", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.authorsUseCaseMock.EXPECT().AddAuthor(ctx, basicAuthor.Name).
			Return(nilAuthor, entity.ErrAuthorAlreadyExists)

		result, err := controllerSetup.service.RegisterAuthor(ctx, &basicAuthorRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.AlreadyExists, entity.ErrAuthorAlreadyExists.Error()))
	})

	t.Run("author invalid name (unsupported symbol)", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.RegisterAuthor(ctx, &library.RegisterAuthorRequest{
			Name: "@",
		})
		require.Nil(t, result)
		require.Error(t, err)
	})

	t.Run("author invalid name (invalid length)", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.RegisterAuthor(ctx, &library.RegisterAuthorRequest{
			Name: "",
		})
		require.Nil(t, result)
		require.Error(t, err)

		result, err = controllerSetup.service.RegisterAuthor(ctx, &library.RegisterAuthorRequest{
			Name: strings.Repeat("a", 1000),
		})
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestController_GetAuthor(t *testing.T) {
	t.Parallel()

	basicAuthor := entity.Author{
		ID:   uuid.Nil,
		Name: "author",
	}

	nilAuthor := entity.Author{}

	basicAuthorRequest := library.GetAuthorInfoRequest{
		Id: basicAuthor.ID.String(),
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.authorsUseCaseMock.EXPECT().GetAuthor(ctx, basicAuthor.ID).
			Return(basicAuthor, nil)

		result, err := controllerSetup.service.GetAuthorInfo(ctx, &basicAuthorRequest)
		require.Equal(t, result, &library.GetAuthorInfoResponse{
			Id:   basicAuthor.ID.String(),
			Name: basicAuthor.Name,
		})
		require.NoError(t, err)
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.authorsUseCaseMock.EXPECT().GetAuthor(ctx, basicAuthor.ID).
			Return(nilAuthor, entity.ErrAuthorNotFound)

		result, err := controllerSetup.service.GetAuthorInfo(ctx, &basicAuthorRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.NotFound, entity.ErrAuthorNotFound.Error()))
	})

	t.Run("author invalid uuid", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.GetAuthorInfo(ctx, &library.GetAuthorInfoRequest{
			Id: "not uuid",
		})
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestController_UpdateAuthor(t *testing.T) {
	t.Parallel()

	basicAuthor := entity.Author{
		ID:   uuid.Nil,
		Name: "author",
	}

	nilAuthor := entity.Author{}

	updatedBasicAuthor := entity.Author{
		ID:   basicAuthor.ID,
		Name: basicAuthor.Name + " update",
	}

	basicAuthorRequest := library.ChangeAuthorInfoRequest{
		Id:   updatedBasicAuthor.ID.String(),
		Name: updatedBasicAuthor.Name,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.authorsUseCaseMock.EXPECT().UpdateAuthor(ctx, updatedBasicAuthor.ID, updatedBasicAuthor.Name).
			Return(updatedBasicAuthor, nil)

		result, err := controllerSetup.service.ChangeAuthorInfo(ctx, &basicAuthorRequest)
		require.NotNil(t, result)
		require.NoError(t, err)
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)
		controllerSetup.authorsUseCaseMock.EXPECT().UpdateAuthor(ctx, updatedBasicAuthor.ID, updatedBasicAuthor.Name).
			Return(nilAuthor, entity.ErrAuthorNotFound)

		result, err := controllerSetup.service.ChangeAuthorInfo(ctx, &basicAuthorRequest)
		require.Nil(t, result)
		require.ErrorIs(t, err, status.Error(codes.NotFound, entity.ErrAuthorNotFound.Error()))
	})

	t.Run("author invalid uuid", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.ChangeAuthorInfo(ctx, &library.ChangeAuthorInfoRequest{
			Id:   "not uuid",
			Name: "author",
		})
		require.Nil(t, result)
		require.Error(t, err)
	})

	t.Run("author invalid name (unsupported symbol)", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.ChangeAuthorInfo(ctx, &library.ChangeAuthorInfoRequest{
			Id:   uuid.Nil.String(),
			Name: "@",
		})
		require.Nil(t, result)
		require.Error(t, err)
	})

	t.Run("author invalid name (unsupported length)", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		t.Cleanup(func() {
			controller.Finish()
		})
		ctx := context.Background()

		controllerSetup := createControllerSetup(controller)

		result, err := controllerSetup.service.ChangeAuthorInfo(ctx, &library.ChangeAuthorInfoRequest{
			Id:   uuid.Nil.String(),
			Name: strings.Repeat("a", 1000),
		})
		require.Nil(t, result)
		require.Error(t, err)
	})
}

type controllerSetup struct {
	booksUseCaseMock   generatedmocks.MockBooksUseCase
	authorsUseCaseMock generatedmocks.MockAuthorsUseCase
	service            library.LibraryServer
}

func createControllerSetup(controller *gomock.Controller) *controllerSetup {
	booksUseCaseMock := generatedmocks.NewMockBooksUseCase(controller)
	authorsUseCaseMock := generatedmocks.NewMockAuthorsUseCase(controller)
	logger := zap.NewNop()

	return &controllerSetup{
		booksUseCaseMock:   *booksUseCaseMock,
		authorsUseCaseMock: *authorsUseCaseMock,
		service:            New(logger, booksUseCaseMock, authorsUseCaseMock),
	}
}
