alter table temp_checker.temperature_data rename column temperature_data_id to sensor_data_id;

alter table temp_checker.temperature_data rename column value to temperature;

alter table temp_checker.temperature_data rename to sensor_data;
