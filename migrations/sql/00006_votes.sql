-- +goose Up
-- +goose StatementBegin
CREATE TABLE votes(
    id SERIAL PRIMARY KEY,
    election_id INT REFERENCES elections(id) ON DELETE CASCADE,
    candidate_id INT REFERENCES candidates(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),

    UNIQUE(user_id, election_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE votes;
-- +goose StatementEnd
