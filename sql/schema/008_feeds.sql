-- +goose Up
ALTER TABLE feeds
ADD CONSTRAINT unique_url UNIQUE (url);

-- +goose Down
ALTER TABLE feeds
DROP CONSTRAINT unique_url;