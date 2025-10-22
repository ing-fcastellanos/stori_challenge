-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    transactions (
        id SERIAL PRIMARY KEY,
        legacy_id INT NOT NULL,
        user_id INT NOT NULL,
        amount DECIMAL(10, 2) NOT NULL,
        timestamp TIMESTAMP NOT NULL
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;

-- +goose StatementEnd