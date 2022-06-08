#!/usr/bin/env sh
systemctl stop corona-dashboard
userdel -f corona-dashboard >/dev/null
exit 0
