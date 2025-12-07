alter table temp_checker.temperature_data rename column sensor_data_id to temperature_data_id;

alter table temp_checker.temperature_data rename column temperature to value;

alter table temp_checker.sensor_data rename to temperature_data;