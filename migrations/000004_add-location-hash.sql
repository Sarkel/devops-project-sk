-- +goose Up
alter table temp_checker.location add column location_sid varchar(10) not null unique;

alter table temp_checker.location_sensor alter column sensor_id type varchar(10) using sensor_id::varchar(10);

alter table temp_checker.location_sensor rename column sensor_id to sensor_sid;

-- +goose Down
alter table temp_checker.location_sensor rename column sensor_sid to sensor_id;

alter table temp_checker.location drop column location_sid;

alter table temp_checker.location_sensor alter column sensor_id type varchar(20) using sensor_id::varchar(20);
