#!/bin/sh
getent group corona-dashboard >/dev/null || \
	groupadd -r corona-dashboard
getent passwd corona-dashboard >/dev/null || \
	useradd -r -g corona-dashboard -s /sbin/nologin \
    -c "User for the Corona-Dashboard Microservice" corona-dashboard
mkdir -p /var/log/corona-dashboard
chmod 0644 /var/log/corona-dashboard
chown corona-dashboard:corona-dashboard /var/log/corona-dashboard
exit 0
