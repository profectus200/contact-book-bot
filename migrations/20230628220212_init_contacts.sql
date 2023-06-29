-- +goose Up
-- +goose StatementBegin

CREATE TABLE contacts
(
    tg_user_id  BIGINT REFERENCES users (tg_user_id),
    contact_id  INTEGER,
    name        TEXT,
    phone       TEXT,
    birthday    DATE,
    description TEXT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE contacts;

-- +goose StatementEnd
