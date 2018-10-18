<?php
// [@[#[!NO-STRIP!]#]@]
// [CFG - SETTINGS / ADMIN]
// v.3.5.7 r.2017.09.05 / smart.framework.v.3.5

//----------------------------------------------------- PREVENT EXECUTION BEFORE RUNTIME READY
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------

//--------------------------------------- Templates and Home Page
$configs['app']['admin-domain'] 					= 'mydomain.ext'; 				// admin domain as yourdomain.ext
$configs['app']['admin-home'] 						= 'cloud.welcome';			// admin home page action
$configs['app']['admin-default-module'] 			= 'cloud';					// admin default module
$configs['app']['admin-template-path'] 				= 'default';				// default admin templates folder from etc/templates/
$configs['app']['admin-template-file'] 				= 'template.htm';			// default admin template file
//---------------------------------------

//--
$configs['app']['url']								= 'https://'.$configs['app']['admin-domain'].'/ncloud/';
//--
//define('APP_INSTALL_PASSWORD', '...');
define('APP_AUTH_PRIVILEGES', '<admin>,<webdav>,<caldav>,<carddav>');
$configs['app-auth']['adm-namespaces'] = [
	'Cloud - Files' 		=> $configs['app']['url'].'admin.php/page/cloud.files/~',
	'Cloud - Calendar' 		=> $configs['app']['url'].'admin.php/page/cloud.ical/~',
	'Cloud - WebCalendar' 		=> $configs['app']['url'].'admin.php?/page/cloud.icalweb',
	'Cloud - Addressbook' 		=> $configs['app']['url'].'admin.php/page/cloud.abook/~',
	'Cloud - WebAddressbook' 	=> $configs['app']['url'].'admin.php?/page/cloud.abookweb'
];
//--
//define('NCLOUD_WEBDAV_PROPFIND_ETAG_MAX_FSIZE', 25000000); // PROPFIND ETag up to 25 MB Files (this will slow down things and is good to be enabled only if sync operations are used ...)
define('NCLOUD_WEBDAV_SHOW_QUOTA', true);  // files
define('NCLOUD_CALDAV_SHOW_QUOTA', true);  // ical
define('NCLOUD_CARDDAV_SHOW_QUOTA', true); // abook
//--

// end of php code
?>