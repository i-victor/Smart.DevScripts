#!/usr/bin/env python2.7

import datetime
import sys

print "Python version is: " , sys.version_info
print "Testing Perl Year 2038 Bug on UNIXW-SERVER";
print "The MAX INTEGER is: " , sys.maxsize
print "the Date on the button should be: 2039-01-01 and it is: " + (datetime.datetime.fromtimestamp(int("2177452800")).strftime('%Y-%m-%d %H:%M:%S'))

#END
