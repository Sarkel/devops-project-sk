-- +goose Up
create index location_sensor_location_id_index
    on temp_checker.location_sensor (location_id);

create index sensor_data_location_sensor_id_index
    on temp_checker.sensor_data (location_sensor_id);

create index sensor_data_timestamp_index
    on temp_checker.sensor_data (timestamp desc);

-- +goose Down
drop index if exists temp_checker.location_sensor_location_id_index;

drop index if exists temp_checker.sensor_data_location_sensor_id_index;

drop index if exists temp_checker.sensor_data_timestamp_index;
