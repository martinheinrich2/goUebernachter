-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE guests
(
    id             INTEGER     NOT NULL PRIMARY KEY,
    last_name      varchar(50) NOT NULL,
    first_name     varchar(50) NOT NULL,
    birth_date     date        NOT NULL,
    birth_place    varchar(50) NOT NULL,
    id_number      varchar(50) NOT NULL,
    nationality    varchar(50) NOT NULL,
    last_residence varchar(50) NOT NULL,
    house_ban      bool        NOT NULL,
    hb_end_date    date        NULL,
    hb_start_date  date        NULL
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS guests;