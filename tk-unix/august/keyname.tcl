#!/usr/local/bin/wish8.5
tk_messageBox -type ok -title "INSTRUCTION" -message "Just press any key and you will get a message which tells you the Tcl/Tk name of the key."
bind all <Key> {tk_messageBox -type ok -title "FIND KEYNAME" -message "Key %K pressed"}