-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    currency_id BIGSERIAL,
    day_limit NUMERIC(10, 5) NOT NULL DEFAULT 0,
    week_limit NUMERIC(10, 5) NOT NULL DEFAULT 0,
    month_limit NUMERIC(10, 5) NOT NULL DEFAULT 0,
    CONSTRAINT day_limit CHECK (day_limit >= 0),
    CONSTRAINT week_limit CHECK (week_limit >= 0),
    CONSTRAINT month_limit CHECK (month_limit >= 0)
);
-- Индекса по умолчанию типа btree достаточно, доступ только по полю id
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
