#!/usr/bin/env bash

haproxy -f /etc/haproxy/haproxy.cfg -f /etc/haproxy/stream.d/ -D

exec /usr/bin/supervisord -c /etc/supervisord.conf