#!/bin/sh

doas ifconfig iwm0 nwid "Some Network" wpakey 0x...(hex key here ...)
doas pfctl -f /etc/pf.conf

