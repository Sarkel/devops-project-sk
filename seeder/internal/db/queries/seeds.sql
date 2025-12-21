-- name: CreateLocation :one
insert into temp_checker.location (location_name, latitude, longitude, location_sid)
values ($1, $2, $3, $4)
on conflict (location_sid) do update set location_sid = excluded.location_sid
returning location_id;

-- name: CreateLocationSensor :one
insert into temp_checker.location_sensor (location_id, sensor_sid, type)
values ($1, $2, $3)
on conflict (location_id, sensor_sid) do nothing
returning location_sensor_id;

-- name: GetAllLocationIDs :many
select location_id from temp_checker.location;

-- name: GetAllLocationSensorIDs :many
select location_sensor_id from temp_checker.location_sensor;

-- name: CreateSensorData :exec
insert into temp_checker.sensor_data (location_sensor_id, temperature, timestamp)
values ($1, $2, $3);

-- name: DeleteAllLocations :exec
delete from temp_checker.location;

-- name: DeleteAllLocationSensors :exec
delete from temp_checker.location_sensor;

-- name: DeleteAllSensorData :exec
delete from temp_checker.sensor_data;

-- name: GetAllLocations :many
select location_id, location_name, latitude, longitude, location_sid
from temp_checker.location
order by location_id;

-- name: GetAllLocationSensors :many
select location_sensor_id, location_id, sensor_sid, type
from temp_checker.location_sensor
order by location_sensor_id;

-- name: GetAllSensorData :many
select sensor_data_id, location_sensor_id, temperature, timestamp
from temp_checker.sensor_data
order by timestamp desc;
