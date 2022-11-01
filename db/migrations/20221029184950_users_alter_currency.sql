-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN currency_id RESTRICT,
    ADD  COLUMN currency varchar(5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN currency,
    ADD  COLUMN currency_id BIGSERIAL;
-- +goose StatementEnd
