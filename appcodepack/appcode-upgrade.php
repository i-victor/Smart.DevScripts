<?php
// AppCodePack - a PHP, JS and CSS Optimizer / NetArchive Pack Upgrade Script
// (c) 2006-2017 unix-world.org - all rights reserved

//----------------------------------------------------- PREVENT EXECUTION BEFORE RUNTIME READY
if(!defined('APPCODEPACK_APP_ID')) { // this must be defined in the first line of the application
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------

//#####
// Sample AppCodePack Upgrade Script, r.171122.1705
// CUSTOMIZE IT AS NEEDED
//#####

//--
// THIS IS A SAMPLE UPGRADE SCRIPT THAT WILL BE RUN AFTER THE PACKAGE DEPLOYMENT
// AS THIS SCRIPT WILL RUN INSIDE THE APPCODEUNPACK (UNDER A DIFFERENT DIRECTORY TAKE IN CONSIDERATION THIS ASPECT ...)
// THIS SCRIPT IS VERY LIMITED AND MUST BE USED FOR SUCCESSFUL AFTER-DEPLOYMENT TASKS LIKE:
// 		* clear Redis cache
// 		* upgrade the SQL databases
//--
// IMPORTANT:
// 		* Below this line, this script should not die() but only throw catcheable exceptions of the \Exception object because will terminate the parent script (appcodeunpack.php) prematurely ...
// 		* This script should not output anything
// 		* If script result is OK, no exception will be throw, thus considered as SUCCESS
//--

//===== Code below is only sample, and can be removed if not needed =====

//-- 1st, clear redis cache DB#17 after deployment
AppCodePackUpgrade::RunCmd('/usr/local/bin/redis-cli -n 17 FLUSHDB'); // throws if unsuccessful
//--

//-- if success, remove maintenance.html (need to exit maintenance mode before run the next command)
AppCodePackUpgrade::RemoveMaintenanceFile(); // throws if unsuccessful
//--

//-- if success, run a post deployment task (by example, check if website works ...)
$arr = (array) AppCodePackUpgrade::RunCmd('/usr/local/bin/curl -s -o /dev/null -w '.escapeshellarg('%{http_code}').' --get --connect-timeout 30 --max-time 150 --insecure --url '.escapeshellarg('https://127.0.0.1'));
//--
if((string)trim((string)$arr['stdout']) != '200') { // expect HTTP Status 200
	throw new Exception('CURL GET - HTTP Status FAILED # Expect 200 but Result is: '.(string)trim((string)$arr['stdout']));
} //end if
//--


//=====================================================================================
//===================================================================================== CLASS START
//=====================================================================================


final class AppCodePackUpgrade {

	// ::
	// v.170921

	// run a command and if not successful throw error
	// returns: - ; Throws Error if not successful
	public static function RunCmd($cmd) {
		//--
		$parr = (array) AppPackUtils::run_proc_cmd(
			(string) $cmd,
			null,
			null,
			null
		);
		$exitcode = $parr['exitcode']; // don't make it INT !!!
		$stdout = (string) $parr['stdout'];
		$stderr = (string) $parr['stderr'];
		if(($exitcode !== 0) OR ((string)$stderr != '')) { // exitcode is zero (0) on success and no stderror
			throw new Exception(__METHOD__.'() :: FAILED to run command ['.$cmd.'] # ExitCode: '.$exitcode.' ; Errors: '.$stderr);
			return (array) $parr;
		} //end if
		//--
		return (array) $parr;
		//--
	} //END FUNCTION


	// clear the maintenance.html file (may be needed if need to run a command after maintenance has been disabled ...)
	// returns: - ; Throws Error if not successful
	public static function RemoveMaintenanceFile() {
		//--
		return (int) self::RemoveAppFile('maintenance.html');
		//--
	} //END FUNCTION


	// remove a file inside app folder (may be needed for some upgrades to remove temporary task files ...)
	// returns: - ; Throws Error if not successful
	public static function RemoveAppFile($file_path) {
		//--
		$file_path = AppPackUtils::safe_pathname((string)$file_path);
		if((string)$file_path == '') {
			throw new Exception(__METHOD__.'() :: Empty FilePath');
			return 0;
		} //end if
		//--
		if((string)APPCODEPACK_APP_ID == '') {
			throw new Exception(__METHOD__.'() :: Empty APPCODEPACK_APP_ID');
			return 0;
		} //end if
		//--
		if(!AppPackUtils::check_if_safe_file_or_dir_name((string)APPCODEPACK_APP_ID)) {
			throw new Exception(__METHOD__.'() :: Unsafe APPCODEPACK_APP_ID: '.APPCODEPACK_APP_ID);
			return 0;
		} //end if
		//--
		$file_app_path = (string) AppPackUtils::add_dir_last_slash((string)APPCODEPACK_APP_ID).$file_path;
		if(!AppPackUtils::check_if_safe_path((string)$file_app_path)) {
			throw new Exception(__METHOD__.'() :: Unsafe Path: '.$file_app_path);
			return 0;
		} //end if
		//--
		if(AppPackUtils::is_type_file((string)$file_app_path)) { // this scripts runs in the parent of {app-id}/
			AppPackUtils::delete((string)$file_app_path);
			if(AppPackUtils::path_exists((string)$file_app_path)) {
				throw new Exception(__METHOD__.'() :: FAILED to remove the file: '.(string)$file_app_path);
				return 0;
			} //end if
		} //end if
		//--
		return 1;
		//--
	} //END FUNCTION


} //END CLASS


//=====================================================================================
//===================================================================================== CLASS END
//=====================================================================================


//end of php code
?>