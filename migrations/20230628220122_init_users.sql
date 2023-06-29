-- +goose Up
-- +goose StatementBegin

CREATE TABLE users
(
    tg_user_id       BIGINT UNIQUE PRIMARY KEY,
    contact_id       INTEGER,
    message_id       INTEGER,
    current_state    TEXT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users;

-- +goose StatementEnd
