ulimit -c 0
export XKL_XMODMAP_DISABLE=1
export GPG_AGENT_INFO=/var/run/user/`id -u`/gnupg/S.gpg-agent
#gpgconf --launch gpg-agent
exec dbus-launch mate-session
#exec dbus-launch startxfce4
