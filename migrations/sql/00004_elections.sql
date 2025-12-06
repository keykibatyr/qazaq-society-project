-- +goose Up
-- +goose StatementBegin
CREATE TABLE elections (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    start_at TIMESTAMP,
    end_at TIMESTAMP,
    published BOOLEAN,
    created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE elections;
-- +goose StatementEnd
