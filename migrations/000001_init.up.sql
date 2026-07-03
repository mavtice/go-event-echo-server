CREATE SCHEMA event_echo;

CREATE TABLE event_echo.users (
    id                 SERIAL       PRIMARY KEY,
    version            BIGINT       NOT NULL DEFAULT 1,

    first_name         VARCHAR(50)  NOT NULL CHECK (char_length(first_name) BETWEEN 3 AND 50),
    last_name          VARCHAR(50)  NOT NULL CHECK (char_length(last_name) BETWEEN 3 AND 50),
    date_of_birthday   DATE         NOT NULL,

    email              VARCHAR(100) NOT NULL CHECK (
        email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'
    ),

    password           VARCHAR(256) NOT NULL,
    role               VARCHAR(5)   NOT NULL,
    notification_token VARCHAR(256)
);
