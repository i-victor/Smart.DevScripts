#!/usr/bin/env php-5.6
<?php

// (c) 2008-2018 unix-world.org

ini_set('display_errors', '1'); // display runtime errors
error_reporting(E_ALL & ~E_NOTICE & ~E_STRICT); // error reporting
date_default_timezone_set('UTC');

echo "Testing PHP Year 2038 Bug on UNIXW-SERVER\n";
echo "The PHP version is: ".phpversion()."\n";
echo "PHP_INT_SIZE on 64-bit should be 8 and is: ".PHP_INT_SIZE."\n";
echo "PHP_INT_MAX is: ".PHP_INT_MAX."\n";
for ($clock = 2147483641; $clock < 2147483651; $clock++) {
    echo date('Y-m-d H:i:s', ($clock))."\n";
} //end for
echo "#END\n";

?>