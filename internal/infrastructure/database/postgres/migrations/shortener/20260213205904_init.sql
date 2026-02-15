-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    short_url    VARCHAR(10)  PRIMARY KEY,
    original_url  TEXT         NOT NULL UNIQUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_urls_original_url ON urls (original_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
DROP INDEX IF EXISTS idx_urls_original_url;
-- +goose StatementEnd