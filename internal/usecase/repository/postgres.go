package repository

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/project/library/internal/entity"
)

var _ BooksRepository = (*postgresRepository)(nil)
var _ AuthorsRepository = (*postgresRepository)(nil)

type postgresRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(db *pgxpool.Pool, logger *zap.Logger) *postgresRepository {
	return &postgresRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *postgresRepository) AddAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	repo.logger.Info("AddAuthor repo: started")
	defer repo.logger.Info("AddAuthor repo: finished")

	return withTransaction(ctx, repo.db, repo.logger, func(tx pgx.Tx) (entity.Author, error) {
		const query = `
		INSERT INTO author (name) 
		VALUES ($1) 
		RETURNING id
	`
		err := tx.QueryRow(ctx, query, author.Name).Scan(&author.ID)
		if err != nil {
			return entity.Author{}, fmt.Errorf("add author query failed: %w", err)
		}

		return author, nil
	})
}

func (repo *postgresRepository) GetAuthor(ctx context.Context, authorID uuid.UUID) (entity.Author, error) {
	repo.logger.Info("GetAuthor repo: started")
	defer repo.logger.Info("GetAuthor repo: finished")
	const query = `
        SELECT id, name
        FROM author 
        WHERE id = $1
    `
	var author entity.Author
	err := repo.db.QueryRow(ctx, query, authorID).Scan(&author.ID, &author.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Author{}, entity.ErrAuthorNotFound
	}

	if err != nil {
		return entity.Author{}, fmt.Errorf("get author query failed: %w", err)
	}

	return author, nil
}

func (repo *postgresRepository) UpdateAuthor(ctx context.Context, authorID uuid.UUID, name string) (entity.Author, error) {
	repo.logger.Info("UpdateAuthor repo: started")
	defer repo.logger.Info("UpdateAuthor repo: finished")
	return withTransaction(ctx, repo.db, repo.logger, func(tx pgx.Tx) (entity.Author, error) {
		const query = `
        UPDATE author
		SET name = $1
		WHERE id = $2
		RETURNING id, name
    `
		var author entity.Author
		err := tx.QueryRow(ctx, query, name, authorID).Scan(&author.ID, &author.Name)
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Author{}, entity.ErrAuthorNotFound
		}

		if err != nil {
			return entity.Author{}, fmt.Errorf("update author query failed: %w", err)
		}

		return author, nil
	})
}

func (repo *postgresRepository) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	repo.logger.Info("AddBook repo: started")
	defer repo.logger.Info("AddBook repo: finished")
	return withTransaction(ctx, repo.db, repo.logger, func(tx pgx.Tx) (entity.Book, error) {
		const queryAddBook = `
		INSERT INTO book (name) 
		VALUES ($1) 
		RETURNING id, name, created_at , updated_at
	`

		var addedBook entity.Book
		err := tx.QueryRow(ctx, queryAddBook, book.Name).Scan(&addedBook.ID, &addedBook.Name, &addedBook.CreatedAt, &addedBook.UpdatedAt)

		if err != nil {
			return entity.Book{}, fmt.Errorf("add book query failed: %w", err)
		}

		if err := insertAuthorsBatch(ctx, tx, addedBook.ID, book.AuthorIDs); err != nil {
			return entity.Book{}, fmt.Errorf("add book's authors query failed: %w", err)
		}

		addedBook.AuthorIDs = book.AuthorIDs
		return addedBook, nil
	})
}

func (repo *postgresRepository) GetBook(ctx context.Context, bookID uuid.UUID) (entity.Book, error) {
	repo.logger.Info("GetBook repo: started")
	defer repo.logger.Info("GetBook repo: finished")
	const query = `
        SELECT 
            id,
            name,
            created_at,
            updated_at,
            ARRAY(
                SELECT author_id 
                FROM author_book 
                WHERE book_id = $1
            ) AS author_ids
        FROM book
        WHERE id = $1
    `

	var book entity.Book
	err := repo.db.QueryRow(ctx, query, bookID).Scan(
		&book.ID,
		&book.Name,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.AuthorIDs,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Book{}, entity.ErrBookNotFound
	}
	if err != nil {
		return entity.Book{}, fmt.Errorf("get book query failed: %w", err)
	}

	return book, nil
}

func (repo *postgresRepository) UpdateBook(ctx context.Context, bookID uuid.UUID, name string, authorIDs uuid.UUIDs) (entity.Book, error) {
	repo.logger.Info("UpdateBook repo: started")
	defer repo.logger.Info("UpdateBook repo: finished")
	return withTransaction(ctx, repo.db, repo.logger, func(tx pgx.Tx) (entity.Book, error) {
		const queryUpdateBook = `
		UPDATE book
		SET name = $1
		WHERE id = $2
		RETURNING id, name, created_at, updated_at;
	`

		var book entity.Book
		err := tx.QueryRow(ctx, queryUpdateBook, name, bookID).Scan(&book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Book{}, entity.ErrBookNotFound
		}
		if err != nil {
			return entity.Book{}, fmt.Errorf("update book query failed: %w", err)
		}

		const queryDeleteBookAuthors = `
		DELETE FROM author_book
		WHERE book_id = $1;
	`

		_, err = tx.Exec(ctx, queryDeleteBookAuthors, bookID)
		if err != nil {
			return entity.Book{}, fmt.Errorf("delete previous book's author query failed: %w", err)
		}

		if err := insertAuthorsBatch(ctx, tx, bookID, authorIDs); err != nil {
			return entity.Book{}, fmt.Errorf("add book's authors query failed: %w", err)
		}

		book.AuthorIDs = authorIDs
		return book, nil
	})
}

func (repo *postgresRepository) GetAuthorBooks(ctx context.Context, authorID uuid.UUID) ([]entity.Book, error) {
	repo.logger.Info("GetAuthorBooks repo: started")
	defer repo.logger.Info("GetAuthorBooks repo: finished")
	const query = `
        SELECT 
            book.id, 
            book.name, 
            book.created_at, 
            book.updated_at,
            ARRAY(
                SELECT author_book.author_id 
                FROM author_book 
                WHERE author_book.book_id = book.id
            ) AS author_ids
        FROM book
        INNER JOIN author_book ON book.id = author_book.book_id
        WHERE author_book.author_id = $1
    `

	rows, err := repo.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query author books: %w", err)
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		var authorIDs []uuid.UUID

		err := rows.Scan(
			&book.ID,
			&book.Name,
			&book.CreatedAt,
			&book.UpdatedAt,
			&authorIDs,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan book row: %w", err)
		}

		book.AuthorIDs = authorIDs
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return books, nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}

func withTransaction[T any](ctx context.Context, db *pgxpool.Pool, logger *zap.Logger, fn func(pgx.Tx) (T, error)) (T, error) {
	logger.Info("Transaction started")
	defer logger.Info("Transaction finished")
	var zero T
	tx, err := db.Begin(ctx)
	if err != nil {
		return zero, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			logger.Error("failed to rollback transaction", zap.Error(err))
		}
		logger.Info("Transaction rolled back")
	}(tx, ctx)

	result, err := fn(tx)
	if err != nil {
		return zero, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Info("Not commited", zap.Error(err))
		return zero, err
	}

	logger.Info("Commited")

	return result, nil
}

func insertAuthorsBatch(ctx context.Context, tx pgx.Tx, bookID uuid.UUID, authorIDs []uuid.UUID) error {
	batch := &pgx.Batch{}
	for _, authorID := range authorIDs {
		batch.Queue(
			"INSERT INTO author_book (author_id, book_id) VALUES ($1, $2)",
			authorID, bookID,
		)
	}

	results := tx.SendBatch(ctx, batch)

	for range authorIDs {
		_, err := results.Exec()

		if isForeignKeyViolation(err) {
			return entity.ErrAuthorNotFound
		}
		if err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("insert batch closing failed: %w", err)
	}

	return nil
}
