-- +goose Up
-- +goose StatementBegin
CREATE TABLE expenses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL,
    category VARCHAR(256) NOT NULL,
    price NUMERIC(10, 5) NOT NULL,
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT price_positive CHECK (price > 0),
    CONSTRAINT category_non_empty CHECK (char_length(category) > 0 AND category = lower(category))
);
-- Операции выполняются по полям user_id, time
-- user_id участвует в операциях равно
-- time участвует в операциях сравнения
-- Эти поля числовые, построим по ним btree индекс
CREATE INDEX expernes_idx ON expenses USING btree (user_id, time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX expernes_idx;
DROP TABLE expenses;
-- +goose StatementEnd
