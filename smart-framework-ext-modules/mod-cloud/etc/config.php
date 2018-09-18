<?php
// [@[#[!NO-STRIP!]#]@]
// [CFG - SETTINGS]
// v.3.5.7 r.2017.09.05 / smart.framework.v.3.5

//----------------------------------------------------- PREVENT EXECUTION BEFORE RUNTIME READY
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------

//--------------------------------------- Info URL
$configs['app']['info-url'] 		= 'netvision.ncloud';					// Info URL: this must be someting like `www . mydomain . net`
//---------------------------------------

//--------------------------------------- SQLite related configuration
$configs['sqlite']['timeout'] 		= 60;									// connection timeout
$configs['sqlite']['slowtime'] 		= 0.0025;								// slow query time (for debugging)
//---------------------------------------

//--------------------------------------- REGIONAL SETTINGS
$configs['regional']['language-id']					= 'en';					// Language `en` | `ro` (must exists as defined)
$configs['regional']['decimal-separator']			= '.';					// decimal separator `.` | `,`
$configs['regional']['thousands-separator']			= ',';					// thousand separator `,` | `.` | ` `
$configs['regional']['calendar-week-start']			= '0';					// 0=start on sunday | 1=start on Monday ; used for both PHP and Javascript
$configs['regional']['calendar-date-format-client'] = 'dd.mm.yy';			// Client Date Format - Javascript (allow only these characters: yy mm dd . - [space])
$configs['regional']['calendar-date-format-server']	= 'd.m.Y';				// Server Date Format - PHP (allow only these characters: Y m d . - [space])
//---------------------------------------
$languages = array('en' => '[EN]');											// default associative array of available languages for this software (do not change without installing new languages support files)
//---------------------------------------

// end of php code
?>