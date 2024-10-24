-- +goose Up
ALTER TABLE feeds
RENAME COLUMN name TO title;

-- +goose Down
ALTER TABLE feeds
RENAME COLUMN title TO name;