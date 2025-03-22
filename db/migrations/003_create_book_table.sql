-- +goose Up
CREATE TABLE book
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_book_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE OR REPLACE TRIGGER trigger_update_book_timestamp
    BEFORE UPDATE
    ON book
    FOR EACH ROW
EXECUTE FUNCTION update_book_timestamp();

-- +goose Down
DROP TABLE book;
DROP FUNCTION update_book_timestamp;