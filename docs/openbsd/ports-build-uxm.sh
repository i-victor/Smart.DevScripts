#!/bin/ksh

### OpenBSD 6.3 Ports Automate Build
## (c) 2016-2018 unix-world.org
## r.180410
###
## Instructions:
#1. Must be run as non-privileged user (ex: admin)
#2. chown -R admin:wheel /usr/ports
#3. *Optional* Link the /usr/ports/bulk ; /usr/ports/distfiles ; /usr/ports/packages ; /usr/ports/plist ; /usr/ports/pobj from outside /usr/ports ex: /_PORTS/
#4. Prepare the /etc/doas.conf to allow user 'admin' to run privileged commands with doas (... see example etc/doas.conf under: unix-configs/openbsd) ; test if works by running: doas pfctl -s rules
#5. Prepare the /etc/mk.conf (ex: can add the -m64 param to build or other changes ... see example etc/mk.conf under: unix-configs/openbsd
###

# Actions: list | test | build
UXM_PORTS_ACTION=build
UXM_LOG_DIR="/_PORTS/uxm-build-logs"
UXM_UNPRIV_USER=admin

### DO NOT MODIFY BELOW CODE !

echo "##### Preparing to [${UXM_PORTS_ACTION}] the OpenBSD 6.3 Ports (Uxm) #####"
sleep 5

# incremental number
UXM_PORT_NUM=0
UXM_MAKE_PORT=-1
UXM_PORT_DIR=""
UXM_BASENAME_PORT=""

# try create build logs dir (only when build)
if [ "$UXM_PORTS_ACTION" == "build" ]; then
	if [ ! -d "$UXM_LOG_DIR" ]; then
		mkdir -p "$UXM_LOG_DIR"
	fi
	if [ ! -d "$UXM_LOG_DIR" ]; then
		echo ""
		echo -n "\033[1;31m"
		echo "**** ERROR: The Log Dir for Builds cannot be created: ${UXM_LOG_DIR} ***"
		echo -n "\033[m"
		echo ""
		exit 1
	fi
fi

# extract all lines except comments or commented lines
for UXM_PORT in `sed 's/[[:space:]]*#.*//;/^[[:space:]]*$/d' ports-build-list-and-hints-uxm.txt`; do
	# if line is non-empty
	if [ ! -z "$UXM_PORT" ]; then
		UXM_PORT_NUM=`expr $UXM_PORT_NUM + 1`
		echo -n "\033[1;34m"
		echo -n "${UXM_PORT_NUM}. OpenBSD-Port: "
		echo -n "\033[m"
		echo -n "/usr/ports/";
		echo -n "\033[1;33m"
		echo -n "${UXM_PORT}"
		echo -n "\033[m"
		# test if directory exist
		UXM_PORT_DIR="/usr/ports/$UXM_PORT"
		if [ ! -d "$UXM_PORT_DIR" ]; then
			if [ "$UXM_PORTS_ACTION" == "test" ]; then
				echo ""
				echo -n "\033[1;31m"
				echo -n " !!! WARNING: Does Path does not Exists !!!";
				echo -n "\033[m"
				echo ""
				exit 2
			fi
			if [ "$UXM_PORTS_ACTION" == "build" ]; then
				echo ""
				echo -n "\033[1;31m"
				echo "**** ERROR: CANNOT BUILD THE PORT: ${UXM_PORT} ***"
				echo -n "\033[m"
				echo ""
				exit 3
			fi
		else
			echo -n "\033[1;32m"
			echo -n " OK ... found ... ";
			echo -n "\033[m"
			if [ "$UXM_PORTS_ACTION" == "build" ]; then
				# replace / with _ to generate unique log names
				UXM_BASENAME_PORT=`echo ${UXM_PORT} | tr '/' '_'`
				echo ""
				echo "*** START BUILDING The Port Package: ${UXM_BASENAME_PORT} from ${UXM_PORT} into ${UXM_PORT_DIR} ***"
				cd "$UXM_PORT_DIR"
				make clean
				make clean=depends
				make clean=flavors
				make clean
				echo "Making Package Bulk: ${UXM_PORT}"
				doas -u $UXM_UNPRIV_USER make package BULK=yes 1>> "${UXM_LOG_DIR}/${UXM_BASENAME_PORT}.log" 2>> "${UXM_LOG_DIR}/${UXM_BASENAME_PORT}.errors.log"
				make clean
				make clean=depends
				echo "*** END BUILDING The Port Package: ${UXM_PORT} ***"
			fi
		fi
		echo ""
	fi
done

# script done
echo ""
echo "*** ... DONE [${UXM_PORTS_ACTION}] ..."
echo ""
exit 0

#END
