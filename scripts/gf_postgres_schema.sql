CREATE TABLE users
(
    acct       VARCHAR(50) UNIQUE NOT NULL,
    pwd        VARCHAR(100),
    fullname   VARCHAR(100),
    created_at TIMESTAMPTZ        NOT NULL
        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ        NOT NULL
        DEFAULT CURRENT_TIMESTAMP
);

