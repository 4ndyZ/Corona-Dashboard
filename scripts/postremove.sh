# Package uninstall
uninstall() {
  userdel -f corona-dashboard >/dev/null || :
  systemctl daemon-reload || :
}

# Package uninstall and purge
purge() {
  rm -drf /etc/corona-dashboard || :
  rm -drf /var/log/corona-dashboard || :
}

# Package upgrade
upgrade() {
  :
}

action="$1"
case "$action" in
  "0" | "remove")
    uninstall
    ;;
  "1" | "upgrade")
    upgrade
    ;;
  "purge")
    uninstall
    purge
    ;;
esac
