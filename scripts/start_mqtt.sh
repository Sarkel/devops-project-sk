#!/bin/sh
set -e

if [ -z "$MQTT_BROKER_USERNAME" ] || [ -z "$MQTT_BROKER_PASSWORD" ]; then
    echo "ERROR: MQTT_BROKER_USERNAME or MQTT_BROKER_PASSWORD not defined in env!"
    exit 1
fi

echo "Generating Mosquitto password file for user: $MQTT_BROKER_USERNAME"

touch /mosquitto/config/passwd

mosquitto_passwd -b /mosquitto/config/passwd "$MQTT_BROKER_USERNAME" "$MQTT_BROKER_PASSWORD"

echo "Starting Mosquitto..."
exec /usr/sbin/mosquitto -c /mosquitto/config/mosquitto.conf