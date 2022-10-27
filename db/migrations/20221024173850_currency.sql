-- +goose Up
-- +goose StatementBegin
CREATE TABLE currencies (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(5) NOT NULL,
    ratio NUMERIC(10, 5) NOT NULL,
    time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT code_non_empty CHECK (char_length(code) > 0 AND code = upper(code)),
    CONSTRAINT ratio_positive CHECK (ratio > 0),
    CONSTRAINT code_uniq UNIQUE(code)
);

-- Для доступа по условию равенства. Значений в таблице ожидается <100, индекс не так важен.
CREATE INDEX currencies_code_idx ON currencies USING HASH (code);

INSERT INTO currencies(code, ratio, time) VALUES ('RUB', '1', '0001-01-01 00:00:00');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX currencies_code_idx;
DROP TABLE currencies;
-- +goose StatementEnd
