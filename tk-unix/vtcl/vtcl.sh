#!/bin/sh

PATH_TO_WISH=/usr/local/bin/wish8.5
VTCL_HOME=/opt/tk-unix/vtcl

export PATH_TO_WISH
export VTCL_HOME

exec ${PATH_TO_WISH} ${VTCL_HOME}/vtcl.tcl $*

#END
