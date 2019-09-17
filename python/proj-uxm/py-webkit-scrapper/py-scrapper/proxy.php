<?php

// Webkit Scrapper Proxy
// (c) 2017-2018 Radu.I
// v.180112.1537

// USAGE @ urlScript: "https://127.0.0.1/sites/py-scrapper/proxy.php?urlproxy=https%3A%2F%2F127.0.0.1%2Fsites%2Fpy-scrapper%2Fhandler.php"

ini_set('display_errors', '1');	// display runtime errors
error_reporting(E_ALL & ~E_NOTICE & ~E_STRICT & ~E_DEPRECATED); // error reporting

$debug = false;

$cpost = '';
$url = '';
if((string)$_GET['urlproxy'] != '') {
	$url = (string) trim((string)$_GET['urlproxy']);
} //end if
if(is_array($_POST)) {
	foreach($_POST as $key => $value) {
		if((string)$key == 'urlproxy') {
			$url = (string) trim((string)$value);
		} else {
			if(is_array($value)) {
				for($i=0; $i<count($value); $i++) {
					$cpost .= urlencode($key).'[]='.rawurlencode($value[$i]).'&';
				} //end for
			} else {
				$cpost .= urlencode($key).'='.rawurlencode($value).'&';
			} //end if else
		} //end if else
	} //end foreach
} //end if

if((string)$url == '') {
	http_response_code(400);
	die('Empty URLProxy');
} //end if

if($debug) {
	file_put_contents('data/get.log', print_r($_GET,1), FILE_APPEND | LOCK_EX);
	file_put_contents('data/post.log', print_r($_POST,1), FILE_APPEND | LOCK_EX);
} //end if

$ch = curl_init((string)$url);
curl_setopt($ch, CURLOPT_POST, 1);
curl_setopt($ch, CURLOPT_POSTFIELDS, (string)$cpost);
curl_setopt($ch, CURLOPT_FOLLOWLOCATION, 0);
curl_setopt($ch, CURLOPT_HEADER, 0);  // DO NOT RETURN HTTP HEADERS
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);  // RETURN THE CONTENTS OF THE CALL
curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, false);
curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
$response = curl_exec($ch);
$info = curl_getinfo($ch);
curl_close($ch);
http_response_code($info['http_code']);
header('Access-Control-Allow-Origin: *');
header('Content-Type: '.$info['content_type']);
echo $response;

if($debug) {
	file_put_contents('data/answer.log', 'Answer # HTTP Code: '.$info['http_code']."\n".'Content-Type: '.$info['content_type']."\n".$response."\n\n", FILE_APPEND | LOCK_EX);
} //end if

// end of php code
?>