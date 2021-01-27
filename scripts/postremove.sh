#!/bin/sh
getent passwd corona-dashboard >/dev/null || \
	userdel -f corona-dashboard
getent group corona-dashboard >/dev/null || \
	groupdel corona-dashboard
exit 0
