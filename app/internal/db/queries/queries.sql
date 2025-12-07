-- name: CreateTemperatureData :many
insert into temp_checker.sensor_data(location_sensor_id, temperature, timestamp)
select unnest(sqlc.arg(location_sensor_ids)::int[]),
       unnest(sqlc.arg(temperatues)::float[]),
       unnest(sqlc.arg(timestamps)::timestamptz[])
returning sensor_data_id;

-- name: GetAPILocationSensors :many
select ls.location_sensor_id,
       ls.sensor_sid,
       l.location_sid,
       l.location_name,
       l.latitude,
       l.longitude,
       l.location_id
from temp_checker.location_sensor as ls
         join temp_checker.location as l on ls.location_id = l.location_id
where ls.type = 'api';

-- name: GetLocationSensorBySensorId :one
select ls.location_sensor_id
from temp_checker.location_sensor as ls
         join temp_checker.location as l on ls.location_id = l.location_id
where ls.sensor_sid = $1
  and l.location_sid = $2;

-- name: LocationExistBySid :one
select count(location_id)
from temp_checker.location
where location_sid = $1;

-- name: GetTodaySensorsSummary :many
select ls.type, now()::date as date, avg(sd.temperature) as avg_temperature
from temp_checker.sensor_data sd
         join temp_checker.location_sensor ls on sd.location_sensor_id = ls.location_sensor_id
         join temp_checker.location l on ls.location_id = l.location_id
where l.location_sid = $1
  and sd.timestamp::date = now()::date
group by ls.type;

-- name: GetLocations :many
select location_sid, location_name
from temp_checker.location;

-- name: GetSensorDataPoints :many
select ls.type,
       case when sqlc.arg(aggregation) is distinct from 'day' then sd.timestamp else sd.timestamp::date end as time_dim,
       round(avg(sd.temperature)::numeric, 1)::float                                                as avg_temperature
from temp_checker.sensor_data sd
         join temp_checker.location_sensor ls on sd.location_sensor_id = ls.location_sensor_id
         join temp_checker.location l on ls.location_id = l.location_id
where l.location_sid = sqlc.arg(location_sid)
  and ls.type = any (sqlc.arg(types)::temp_checker.sensor_type[])
  and sd.timestamp between sqlc.arg(start_datetime)::timestamp and sqlc.arg(end_datetime)::timestamp
group by ls.type, time_dim;