-- +goose Up
-- +goose StatementBegin
ALTER Table currencies ALTER ratio SET DATA TYPE NUMERIC(20, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER Table currencies ALTER ratio SET DATA TYPE NUMERIC(10, 5);
-- +goose StatementEnd
