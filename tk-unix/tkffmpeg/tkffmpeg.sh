#!/bin/sh

export LANG=en_US

if [ -e /opt/tk-unix/tkffmpeg/tkffmpeg.tcl ] 
then
	exec /opt/tk-unix/tkffmpeg/tkffmpeg.tcl "$@"
else
	echo "--------------------------------------------------------------------------------------------"
	echo "|          Error!!!You do not have tkffmpeg :), try to find and run  tkffmpeg.tcl          |"
	echo "--------------------------------------------------------------------------------------------"
fi

#END
