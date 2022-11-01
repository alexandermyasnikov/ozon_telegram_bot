-- +goose Up
-- +goose StatementBegin
ALTER TABLE expenses
    DROP CONSTRAINT category_non_empty,
    ADD  CONSTRAINT category_non_empty CHECK (char_length(category) > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE expenses SET category = lower(category);
ALTER TABLE expenses
    DROP CONSTRAINT category_non_empty,
    ADD  CONSTRAINT category_non_empty CHECK (char_length(category) > 0 AND category = lower(category));
-- +goose StatementEnd
