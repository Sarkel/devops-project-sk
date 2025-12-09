#!/bin/sh
set -e

if [ -z "$AUTH_KEY_NAME" ] || [ -z "$AUTH_KEY_VAL" ]; then
    echo "ERROR: AUTH_KEY_NAME or AUTH_KEY_VAL not defined in env!"
    exit 1
fi

if [ -z "$BASIC_AUTH_USER" ] || [ -z "$BASIC_AUTH_PASSWORD" ]; then
    echo "ERROR: BASIC_AUTH_USER or BASIC_AUTH_PASSWORD not defined in env!"
    exit 1
fi

echo "Generating nginx basic auth file for user: $BASIC_AUTH_USER"

mkdir -p /etc/nginx

touch /etc/nginx/.htpasswd
htpasswd -bc /etc/nginx/.htpasswd "$BASIC_AUTH_USER" "$BASIC_AUTH_PASSWORD"

echo "Rendering nginx config with env variables..."
envsubst '${AUTH_KEY_NAME} ${AUTH_KEY_VAL}' \
  < /etc/nginx/conf.d/default.conf \
  > /etc/nginx/conf.d/default.conf.rendered

mv /etc/nginx/conf.d/default.conf.rendered /etc/nginx/conf.d/default.conf

echo "Starting nginx..."
exec nginx -g 'daemon off;'