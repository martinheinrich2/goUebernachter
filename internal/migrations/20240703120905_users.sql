-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE users
(
    id              INTEGER     NOT NULL PRIMARY KEY,
    last_name       VARCHAR(50) NOT NULL,
    first_name      VARCHAR(50) NOT NULL,
    email           VARCHAR(50) NOT NULL,
    job_title       VARCHAR(50) NOT NULL,
    room            VARCHAR(10) NOT NULL,
    admin           INTEGER DEFAULT 0,
    hashed_password CHAR(60)    NOT NULL,
    created         DATETIME    NOT NULL,
    CONSTRAINT staff_uc_email UNIQUE (email)
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS users;