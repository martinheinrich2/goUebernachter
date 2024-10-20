-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS password_reset_tokens
(
    user_id      VARCHAR(36)  NOT NULL,
    token        VARCHAR(128) NOT NULL UNIQUE,
    token_expiry DATETIME     NOT NULL,
    PRIMARY KEY (user_id, token)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS password_reset_tokens;