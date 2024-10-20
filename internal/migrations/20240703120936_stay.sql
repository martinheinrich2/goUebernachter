-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE stay
(
    id               integer     NOT NULL PRIMARY KEY,
    start_date       date        NOT NULL,
    end_date         date        NULL,
    type_of_stay     varchar(50) NOT NULL,
    room             varchar(20) NOT NULL,
    guest_id         bigint      NOT NULL REFERENCES guests (id) DEFERRABLE INITIALLY DEFERRED,
    social_worker_id bigint      NOT NULL REFERENCES users (id) DEFERRABLE INITIALLY DEFERRED,
    user_id          bigint      NOT NULL REFERENCES users (id) DEFERRABLE INITIALLY DEFERRED,
    appointment      date        NULL,
    appointment_done INTEGER     DEFAULT 0,
    stay_processed   INTEGER     DEFAULT 0
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS stay;