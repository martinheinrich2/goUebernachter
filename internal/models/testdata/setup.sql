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
    appointment_done INTEGER DEFAULT 0,
    stay_processed   INTEGER DEFAULT 0
);


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
    active          BOOL    DEFAULT 1,
    CONSTRAINT staff_uc_email UNIQUE (email)
);

INSERT INTO users (last_name, first_name, email, job_title, room, admin, hashed_password, created, active)
VALUES ('Jones',
        'Alice',
        'alice@example.com',
        'Sozialarbeiter',
        'H321',
        0,
        '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
        '2022-01-01 09:18:24'), 1
