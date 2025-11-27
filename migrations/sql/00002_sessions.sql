-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions(
    id SERIAL PRIMARY KEY,
    user_id int REFERNCES users (id) ON DELETE CASCADE,
    token_hash TEST UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
