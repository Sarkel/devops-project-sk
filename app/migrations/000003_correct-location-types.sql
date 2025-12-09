-- +goose Up
alter table temp_checker.location alter column longitude type float using longitude::float;
alter table temp_checker.location alter column latitude type float using latitude::float;

-- +goose Down
alter table temp_checker.location alter column longitude type numeric using longitude::numeric;
alter table temp_checker.location alter column latitude type numeric using latitude::numeric;
