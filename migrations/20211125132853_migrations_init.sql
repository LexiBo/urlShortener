-- +goose Up
CREATE TABLE urls (
                      full_url varchar NOT NULL,
                      short_url varchar NOT NULL
);

-- +goose Down
DROP TABLE urls;

