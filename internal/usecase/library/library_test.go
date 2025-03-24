package library

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/project/library/internal/entity"
	mocks "github.com/project/library/mocks-generated"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestLibrary_BooksUseCase(t *testing.T) {
	t.Parallel()

	basicBook := entity.Book{
		ID:        uuid.Nil,
		Name:      "book",
		AuthorIDs: []uuid.UUID{uuid.Nil},
	}

	nilBook := entity.Book{}

	t.Run("AddBook", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			expectedBook   entity.Book
			expectedError  error
			inputName      string
			inputAuthorIDs []uuid.UUID
		}{
			{
				name:           "success",
				expectedBook:   basicBook,
				expectedError:  nil,
				inputName:      basicBook.Name,
				inputAuthorIDs: basicBook.AuthorIDs,
			},
			{
				name:           "author not found",
				expectedBook:   nilBook,
				expectedError:  entity.ErrAuthorNotFound,
				inputName:      basicBook.Name,
				inputAuthorIDs: basicBook.AuthorIDs,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				controller := gomock.NewController(t)
				t.Cleanup(func() {
					controller.Finish()
				})
				ctx := context.Background()

				testSetup := createSetup(controller)
				testSetup.booksRepositoryMock.EXPECT().AddBook(ctx, gomock.Any()).
					Return(test.expectedBook, test.expectedError)

				result, err := testSetup.booksUseCase.AddBook(ctx, test.inputName, test.inputAuthorIDs)
				require.Equal(t, result, test.expectedBook)
				require.ErrorIs(t, err, test.expectedError)
			})
		}
	})

	t.Run("GetBook", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name          string
			book          entity.Book
			expectedBook  entity.Book
			expectedError error
		}{
			{
				name:          "success",
				book:          basicBook,
				expectedBook:  basicBook,
				expectedError: nil,
			},
			{
				name:          "book bot found",
				book:          basicBook,
				expectedBook:  nilBook,
				expectedError: entity.ErrBookNotFound,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				controller := gomock.NewController(t)
				t.Cleanup(func() {
					controller.Finish()
				})
				ctx := context.Background()

				testSetup := createSetup(controller)
				testSetup.booksRepositoryMock.EXPECT().GetBook(ctx, gomock.Any()).
					Return(test.expectedBook, test.expectedError)

				result, err := testSetup.booksUseCase.GetBook(ctx, test.book.ID)
				require.Equal(t, result, test.expectedBook)
				require.ErrorIs(t, err, test.expectedError)
			})
		}
	})

	t.Run("UpdateBook", func(t *testing.T) {
		t.Parallel()

		updatedBook := entity.Book{
			ID:        basicBook.ID,
			Name:      basicBook.Name + " update",
			AuthorIDs: []uuid.UUID{uuid.Nil, uuid.Max},
		}

		tests := []struct {
			name          string
			book          entity.Book
			expectedBook  entity.Book
			expectedError error
		}{
			{
				name:          "success",
				book:          updatedBook,
				expectedBook:  updatedBook,
				expectedError: nil,
			},
			{
				name:          "book bot found",
				book:          updatedBook,
				expectedBook:  nilBook,
				expectedError: entity.ErrBookNotFound,
			},
			{
				name:          "author not found",
				book:          updatedBook,
				expectedBook:  nilBook,
				expectedError: entity.ErrAuthorNotFound,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				controller := gomock.NewController(t)
				t.Cleanup(func() {
					controller.Finish()
				})
				ctx := context.Background()

				testSetup := createSetup(controller)
				testSetup.booksRepositoryMock.EXPECT().UpdateBook(ctx, test.book.ID, test.book.Name, test.book.AuthorIDs).
					Return(test.expectedBook, test.expectedError)

				result, err := testSetup.booksUseCase.UpdateBook(ctx, test.book.ID, test.book.Name, test.book.AuthorIDs)
				require.Equal(t, result, test.expectedBook)
				require.ErrorIs(t, err, test.expectedError)
			})
		}
	})

	t.Run("GetAuthorBooks", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)
		defer controller.Finish()

		vasya := uuid.Nil

		book1 := entity.Book{ID: uuid.Nil, Name: "Book 1", AuthorIDs: []uuid.UUID{vasya}}
		book2 := entity.Book{ID: uuid.Max, Name: "Book 2", AuthorIDs: []uuid.UUID{vasya}}

		t.Run("success", func(t *testing.T) {
			ctx := context.Background()
			testSetup := createSetup(controller)
			testSetup.booksRepositoryMock.EXPECT().
				GetAuthorBooks(ctx, vasya).
				Return([]entity.Book{book1, book2}, nil)

			result, err := testSetup.booksUseCase.GetAuthorBooks(context.Background(), vasya)
			require.NoError(t, err)
			require.ElementsMatch(t, []entity.Book{book1, book2}, result)
		})

		t.Run("author not found", func(t *testing.T) {
			ctx := context.Background()
			testSetup := createSetup(controller)
			testSetup.booksRepositoryMock.EXPECT().
				GetAuthorBooks(ctx, vasya).
				Return(nil, entity.ErrAuthorNotFound)

			result, err := testSetup.booksUseCase.GetAuthorBooks(context.Background(), vasya)
			require.ErrorIs(t, entity.ErrAuthorNotFound, err)
			require.Nil(t, result)
		})
	})
}

func TestLibrary_AuthorsUseCase(t *testing.T) {
	t.Parallel()

	basicAuthor := entity.Author{
		ID:   uuid.Nil,
		Name: "author",
	}

	nilAuthor := entity.Author{}

	t.Run("AddAuthor", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			author         entity.Author
			expectedAuthor entity.Author
			expectedError  error
		}{
			{
				name:           "success",
				author:         basicAuthor,
				expectedAuthor: basicAuthor,
				expectedError:  nil,
			},
			{
				name:           "author already exists",
				author:         basicAuthor,
				expectedAuthor: nilAuthor,
				expectedError:  entity.ErrAuthorAlreadyExists,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				controller := gomock.NewController(t)
				t.Cleanup(func() {
					controller.Finish()
				})
				ctx := context.Background()

				testSetup := createSetup(controller)
				testSetup.authorsRepositoryMock.EXPECT().AddAuthor(ctx, gomock.Any()).
					Return(test.expectedAuthor, test.expectedError)

				result, err := testSetup.authorsUseCase.AddAuthor(ctx, test.author.Name)
				require.Equal(t, result, test.expectedAuthor)
				require.ErrorIs(t, err, test.expectedError)
			})
		}
	})

	t.Run("GetAuthor", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			author         entity.Author
			expectedAuthor entity.Author
			expectedError  error
		}{
			{
				name:           "success",
				author:         basicAuthor,
				expectedAuthor: basicAuthor,
				expectedError:  nil,
			},
			{
				name:           "author not found",
				author:         basicAuthor,
				expectedAuthor: nilAuthor,
				expectedError:  entity.ErrAuthorNotFound,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				controller := gomock.NewController(t)
				t.Cleanup(func() {
					controller.Finish()
				})
				ctx := context.Background()

				testSetup := createSetup(controller)
				testSetup.authorsRepositoryMock.EXPECT().GetAuthor(ctx, test.author.ID).
					Return(test.expectedAuthor, test.expectedError)

				result, err := testSetup.authorsUseCase.GetAuthor(ctx, test.author.ID)
				require.Equal(t, result, test.expectedAuthor)
				require.ErrorIs(t, err, test.expectedError)
			})
		}
	})

	t.Run("UpdateAuthor", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			author         entity.Author
			expectedAuthor entity.Author
			expectedError  error
		}{
			{
				name:   "success",
				author: basicAuthor,
				expectedAuthor: entity.Author{
					ID:   basicAuthor.ID,
					Name: basicAuthor.Name + " update",
				},
				expectedError: nil,
			},
			{
				name:          "author not found",
				author:        basicAuthor,
				expectedError: entity.ErrAuthorNotFound,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				controller := gomock.NewController(t)
				t.Cleanup(func() {
					controller.Finish()
				})
				ctx := context.Background()

				testSetup := createSetup(controller)
				testSetup.authorsRepositoryMock.EXPECT().UpdateAuthor(ctx, test.author.ID, test.author.Name).
					Return(test.expectedAuthor, test.expectedError)

				author, err := testSetup.authorsUseCase.UpdateAuthor(ctx, test.author.ID, test.author.Name)
				require.Equal(t, author, test.expectedAuthor)
				require.ErrorIs(t, err, test.expectedError)
			})
		}
	})
}

type setup struct {
	booksRepositoryMock   *mocks.MockBooksRepository
	authorsRepositoryMock *mocks.MockAuthorsRepository
	booksUseCase          BooksUseCase
	authorsUseCase        AuthorsUseCase
}

func createSetup(controller *gomock.Controller) setup {
	authorsRepository := mocks.NewMockAuthorsRepository(controller)
	booksRepository := mocks.NewMockBooksRepository(controller)
	logger := zap.NewNop()
	lib := New(logger, booksRepository, authorsRepository)

	return setup{
		booksRepositoryMock:   booksRepository,
		authorsRepositoryMock: authorsRepository,
		booksUseCase:          lib,
		authorsUseCase:        lib,
	}
}
