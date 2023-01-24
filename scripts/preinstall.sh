# Initial package installation
install() {
  # Create the user and group
  getent group corona-dashboard >/dev/null || groupadd -r corona-dashboard || :
  getent passwd corona-dashboard >/dev/null || useradd -r -g corona-dashboard -s /sbin/nologin \
     -c "User for the Corona-Dashboard Microservice" corona-dashboard || :
}

# Package upgrade
upgrade() {
  :
}

action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  # deb passes $1=configure
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
  # deb passes $1=configure $2=<current version>
  action="upgrade"
fi

case "$action" in
  "1" | "install")
    install
    ;;
  "2" | "upgrade")
    upgrade
    ;;
esac