-- +goose Up
-- +goose StatementBegin
CREATE TABLE candidates(
    id SERIAL PRIMARY KEY,
    election_id INT REFERENCES elections(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    image_url TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE candidates;
-- +goose StatementEnd
