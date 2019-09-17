#!/bin/sh

cd /root
chown -R admin:wheel /usr/ports
cd /usr
THE_DATE=`date +%Y%m%d_%H%M%S`
tar -czf ports-openbsd63-uxm-${THE_DATE}.tar.gz ./ports/
md5 ports-openbsd63-uxm-${THE_DATE}.tar.gz > ports-openbsd63-uxm-${THE_DATE}.tar.gz.checksum
sha256 ports-openbsd63-uxm-${THE_DATE}.tar.gz >> ports-openbsd63-uxm-${THE_DATE}.tar.gz.checksum
sha512 ports-openbsd63-uxm-${THE_DATE}.tar.gz >> ports-openbsd63-uxm-${THE_DATE}.tar.gz.checksum

#END
