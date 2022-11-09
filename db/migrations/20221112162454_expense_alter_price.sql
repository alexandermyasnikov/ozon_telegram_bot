-- +goose Up
-- +goose StatementBegin
ALTER Table expenses ALTER price SET DATA TYPE NUMERIC(20, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER Table expenses ALTER price SET DATA TYPE NUMERIC(10, 5);
-- +goose StatementEnd
