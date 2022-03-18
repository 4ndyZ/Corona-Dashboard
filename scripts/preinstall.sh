#!/usr/bin/env sh
getent group corona-dashboard >/dev/null || \
	groupadd -r corona-dashboard
getent passwd corona-dashboard >/dev/null || \
	useradd -r -g corona-dashboard -s /sbin/nologin \
    -c "User for the Corona-Dashboard Microservice" corona-dashboard
exit 0
