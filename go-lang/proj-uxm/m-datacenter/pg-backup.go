
// GO Lang :: M-DataCenter :: Pg-Backup
// (c) 2020 unix-world.org
// STABLE

package main


import (
	"os"
	"log"
	"fmt"
	"time"

	"github.com/fatih/color"
	smart "github.com/unix-world/smartgo"
)


const (
	PROGR_VERSION = "r.20200510.2259"

	CMD_TIMEOUT = 3600 // 1h per cmd

	PG_HOST = "127.0.0.1"
	PG_PORT = "5432"
	PG_USER = "pgsql"
	PG_PASS = "pgsql"
	PG_DB   = "smart_framework"

	PG_DUMP_FORMAT = "p" // "p" is plain ; "t" is tar

	// (all relative paths) for all below the backup dir must be the same ; !!! the backup dir will be deleted and re-created on each new backup action !!!
	BKP_SCHEMA_FILE = "data-backup/db-schema.sql"
	BKP_DATA_FILE   = "data-backup/db-data.sql"
	BKP_SAFETY_FILE = "data-backup/db-pgdump.ok"

	BKP_ARCHIVE_FOLDER = "data-archive"
)


//---


func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION


//---


func checkIfSafePgDumpConnectionParams() (areValid bool, errDet string) {
	//--
	if(PG_HOST == "") {
		return false, "Empty PgSQL HOST"
	} //end if
	if((!smart.IsNetValidIpAddr(PG_HOST)) && (!smart.IsNetValidHostName(PG_HOST))) {
		return false, "Invalid PgSQL HOST: " + PG_HOST
	} //end if
	if(PG_PORT == "") {
		return false, "Empty PgSQL PORT"
	} //end if
	if(!smart.IsNetValidPortStr(PG_PORT)) {
		return false, "Invalid PgSQL PORT: " + PG_PORT
	} //end if
	//--
	if(PG_USER == "") {
		return false, "Empty PgSQL USER"
	} //end if
	if(PG_PASS == "") {
		return false, "Empty PgSQL PASS"
	} //end if
	//--
	if(PG_DB == "") {
		return false, "Empty PgSQL DB"
	} //end if
	if((PG_DB == "pgsql") || (PG_DB == "postgres") || (PG_DB == "template0") || (PG_DB == "template1")) {
		return false, "Empty or Invalid or Dissalowed PgSQL DB: " + PG_DB
	} //end if
	//--
	if((PG_DUMP_FORMAT != "p") && (PG_DUMP_FORMAT != "t")) {
		return false, "Empty or Invalid PgSQL Dump Format: " + PG_DUMP_FORMAT
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


//---


func checkifSafeBackupDir(fDirName string, checkIfExists bool) (isValid bool, errMsg string) {
	//--
	if((smart.StrTrimWhitespaces(fDirName) == "") ||
		smart.PathIsEmptyOrRoot(fDirName) ||
		smart.PathIsBackwardUnsafe(fDirName) ||
		smart.PathIsAbsolute(fDirName) ||
		smart.StrContains(fDirName, ".") ||
		!smart.StrRegexMatchString(smart.REGEX_SMART_SAFE_PATH_NAME, fDirName)) {
		//--
		return false, "Invalid Backup Dir: " + fDirName
		//--
	} //end if
	//--
	if(checkIfExists == true) {
		if(!smart.PathExists(fDirName) || !smart.PathIsDir(fDirName)) {
			return false, "Non-Existing Backup Dir: " + fDirName
		} //end if
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


func checkIfSafeBackupFilePath(bkpFile string) (isValid bool, errMsg string) {
	//--
	if((smart.StrTrimWhitespaces(bkpFile) == "") ||
		smart.PathIsEmptyOrRoot(bkpFile) ||
		smart.PathIsBackwardUnsafe(bkpFile) ||
		smart.PathIsAbsolute(bkpFile) ||
		smart.StrContains(bkpFile, " ") ||
		!smart.StrRegexMatchString(smart.REGEX_SMART_SAFE_PATH_NAME, bkpFile)) {
		//--
		return false, "Backup File Path (must be relative, non-empty, must not contain spaces, can contain only [ _ a-z A-Z 0-9 - . @ # / ]): " + bkpFile
		//--
	} //end if
	//--
	var fDirName string = smart.PathDirName(bkpFile)
	isDirValid, errDirMsg := checkifSafeBackupDir(fDirName, true)
	if((isDirValid != true) || (errDirMsg != "")) {
		return false, "Invalid Backup Dir: `" + fDirName + "` # " + errDirMsg
	} //end if
	//--
	var fBaseName string = smart.PathBaseName(bkpFile)
	//--
	var theValidFExt = ".sql" // "p", plain
	if(PG_DUMP_FORMAT == "t") {
		theValidFExt = ".tar" // "t", tar
	} //end if else
	//--
	if((smart.StrTrimWhitespaces(fBaseName) == "") ||
		!smart.StrEndsWith(fBaseName, theValidFExt) || (smart.PathBaseExtension(fBaseName) != theValidFExt) ||
		smart.StrStartsWith(fBaseName, ".") ||
		(len(fBaseName) < 5) || (len(smart.StrTrimWhitespaces(fBaseName)) < 5) ||
		!smart.StrRegexMatchString(smart.REGEX_SMART_SAFE_FILE_NAME, fBaseName)) {
		//--
		return false, "Invalid Backup File Name (must ends in .sql, must not start with . and must have min 5 characters): " + fBaseName
		//--
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


//---


func cleanupOldFile(theFile string) (isValid bool, errMsg string) {
	//--
	if(smart.PathExists(theFile)) {
		if(smart.PathIsDir(theFile)) {
			return false, "The cleanup old file path is a Dir: `" + theFile + "`"
		} //end if
		if(smart.PathIsFile(theFile)) {
			smart.SafePathFileDelete(theFile, false)
		} //end if
		if(smart.PathExists(theFile)) {
			return false, "The cleanup old file path cannot be removed: `" + theFile + "`"
		} //end if
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


//---


func getSzFileName(bkpFile string) string {
	//--
	var theSzFile string = bkpFile + ".size"
	//--
	return theSzFile
	//--
} //END FUNCTION


func removeSzFile(bkpFile string) (isOk bool, errMsg string) {
	//--
	var theSzFile string = getSzFileName(bkpFile)
	isClean, errClean := cleanupOldFile(theSzFile)
	if((isClean != true) || (errClean != "")) {
		return false, "Cannot remove the Old Size File: `" + errClean + "`"
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


func createSzFile(bkpFile string) (theSize int64, errMsg string) {
	//--
	fSize, errMsg := smart.SafePathFileGetSize(bkpFile, false)
	if(errMsg != "") {
		return 0, "Cannot Get Size of the Backup File (after backup): `" + bkpFile + "` # " + errMsg
	} //end if
	if(fSize <= 0) {
		return 0, "Backup File is empty (after backup) Size: " + smart.ConvertInt64ToStr(fSize) + " bytes # `" + bkpFile + "`"
	} //end if
	//--
	var theSzFile string = getSzFileName(bkpFile)
	//--
	isSuccess, errWrMsg := smart.SafePathFileWrite(smart.ConvertInt64ToStr(fSize) + "\n", "w", theSzFile, false)
	if((isSuccess != true) || (errWrMsg != "")) {
		return fSize, "Failed to save the Size (" + smart.ConvertInt64ToStr(fSize) + ") of File: `" + bkpFile + "` to `" + theSzFile + "` # " + errWrMsg
	} //end if
	//--
	return fSize, ""
	//--
} //END FUNCTION


//---


func getMd5FileName(bkpFile string) string {
	//--
	var theMd5ChskumFile string = bkpFile + ".md5"
	//--
	return theMd5ChskumFile
	//--
} //END FUNCTION


func removeMd5Checksum(bkpFile string) (isOk bool, errMsg string) {
	//--
	var theMd5ChskumFile string = getMd5FileName(bkpFile)
	isClean, errClean := cleanupOldFile(theMd5ChskumFile)
	if((isClean != true) || (errClean != "")) {
		return false, "Cannot remove the Old MD5 Checksum File: `" + errClean + "`"
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


func createMd5Checksum(bkpFile string) (md5Hash string, errMsg string) {
	//--
	md5Sum, errMsg := smart.SafePathFileMd5(bkpFile, false)
	if((errMsg != "") || ((md5Sum == "") || (len(md5Sum) != 32))) {
		return "", "Failed to get the MD5 sum of File: `" + bkpFile + "` # " + errMsg
	} //end if
	//--
	var theChecksumFileBaseName string = bkpFile
	if(smart.StrContains(theChecksumFileBaseName, "/")) {
		theChecksumFileBaseName = smart.PathBaseName(theChecksumFileBaseName)
	} //end if
	//--
	var theMd5ChskumFile string = getMd5FileName(bkpFile)
	//--
	isSuccess, errWrMsg := smart.SafePathFileWrite("MD5 (" + theChecksumFileBaseName + ") = " + md5Sum + "\n", "w", theMd5ChskumFile, false)
	if((isSuccess != true) || (errWrMsg != "")) {
		return "", "Failed to save the MD5 sum (" + md5Sum + ") of File: `" + bkpFile + "` to `" + theMd5ChskumFile + "` # " + errWrMsg
	} //end if
	//--
	return md5Sum, ""
	//--
} //END FUNCTION


//---


func doBackup(bkpSchemaOrData string) (isBkpValid bool, errBkpMsg string) {
	//-- defs
	var bkpMode string = ""
	var bkpFile string = ""
	var bkpParam string = ""
	if(bkpSchemaOrData == "schema") {
		bkpMode = "SCHEMA"
		bkpFile = BKP_SCHEMA_FILE
		bkpParam = "--schema-only"
	} else if(bkpSchemaOrData == "data") {
		bkpMode = "DATA"
		bkpFile = BKP_DATA_FILE
		bkpParam = "--data-only"
	} else {
		return false, "Invalid DB Backup Mode: `" + bkpSchemaOrData + "`"
	} //end if else
	//-- check file path to be valid
	isValid, errMsg := checkIfSafeBackupFilePath(bkpFile)
	if((isValid != true) || (errMsg != "")) {
		return false, "Invalid DB Backup File Path: `" + errMsg + "`"
	} //end if
	//-- cleanup old backup
	isClean, errClean := cleanupOldFile(bkpFile)
	if((isClean != true) || (errClean != "")) {
		return false, "Cannot cleanup the old backup: `" + errClean + "`"
	} //end if
	//-- cleanup old md5 cksum
	isMd5Clean, errMd5Clean := removeMd5Checksum(bkpFile)
	if((isMd5Clean != true) || (errMd5Clean != "")) {
		return false, "Cannot cleanup the old backup md5: `" + errMd5Clean + "`"
	} //end if
	//-- cleanup old size reg
	isSzClean, errSzClean := removeSzFile(bkpFile)
	if((isSzClean != true) || (errSzClean != "")) {
		return false, "Cannot cleanup the old backup size registration: `" + errSzClean + "`"
	} //end if
	//-- dump
	var pgDumpDetails string = "Host=" + PG_HOST + ":" + PG_PORT + " ; User=" + PG_USER + " ; Pass=***** ; DB=" + PG_DB + " ; File=" + bkpFile
	fmt.Println("========== PgDump " + bkpMode + ": START ==========")
	fmt.Println("PgDump: " + bkpMode + " # " + pgDumpDetails)
	isSuccess, _, errStd := smart.ExecTimedCmd(CMD_TIMEOUT, "output", "capture+output", "PGPASSWORD=" + PG_PASS, "", "pg_dump", "--encoding=UTF8", "--column-inserts", "--blobs", bkpParam, "--no-owner" , "--no-privileges", "--host=" + PG_HOST, "--port=" + PG_PORT, "--user=" + PG_USER, "--format=" + PG_DUMP_FORMAT, "--file=" + bkpFile, PG_DB)
	if((isSuccess != true) || (errStd != "")) {
		return false, "PgDump returned Errors / StdErr:\n`" + errStd + "`\n"
	} //end if
	//-- check that backup file exists, after backup
	if((!smart.PathExists(bkpFile)) || (!smart.PathIsFile(bkpFile))) {
		return false, "Backup File cannot be found (after backup) at: `" + bkpFile + "`"
	} //end if
	//-- create size reg file + check file size to be > 0
	fSize, errSize := createSzFile(bkpFile)
	if((fSize <= 0) || (errSize != "")) {
		return false, "Failed to Create the Size Registration File of Backup File (after backup): `" + bkpFile + "`, Size=" + smart.ConvertInt64ToStr(fSize) + " # " + errSize
	} //end if
	//-- create md5 sum file + check sum to be valid
	md5Hash, errMd5Hash := createMd5Checksum(bkpFile)
	if((md5Hash == "") || (len(md5Hash) != 32) || (errMd5Hash != "")) {
		return false, "Failed to Create the MD5 Sum of Backup File (after backup): `" + bkpFile + "` # " + errMd5Hash
	} //end if
	//-- end messages
	fmt.Println(color.GreenString("OK: DB " + bkpMode + " backup") + " # SIZE(" + smart.ConvertInt64ToStr(fSize) + ") / MD5(" + md5Hash + ") was saved to:", bkpFile)
	fmt.Println("========== PgDump " + bkpMode + ": END ==========")
	//-- @ret
	return true, ""
	//--
} //END FUNCTION


//---


func reCreateDir(theDir string) (isOk bool, errMsg string) {
	//--
	if(smart.PathExists(theDir)) {
		if(smart.PathIsFile(theDir)) {
			smart.SafePathFileDelete(theDir, false)
		} else if(smart.PathIsDir(theDir)) {
			smart.SafePathDirDelete(theDir, false)
		} //end if
		if(smart.PathExists(theDir)) {
			return false, "Cannot Cleanup Old Dir: `" + theDir + "`"
		} //end if
	} //end if
	smart.SafePathDirCreate(theDir, true, false) // allow recursive, deny absolute
	if((!smart.PathExists(theDir)) || (!smart.PathIsDir(theDir))) {
		return false, "Cannot Create the New Dir: `" + theDir + "`"
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


//---


var uxmIsBackupRunning bool = false


func logBackupError(logMessages ...interface{}) {
	//--
	log.Println("[ERROR] :: BACKUP Task :: ", logMessages) // standard logger
	// must not exit
	//--
} //END FUNCTION


func runBackupTask(taskNum int64) {

	//--
	uxmIsBackupRunning = true
	//--

	//--
	var DateTimeStartUtc string = smart.DateNowUtc()
	//--

	//--
	fmt.Println("")
	fmt.Println(color.BlueString("##### M-DataCenter :: PostgreSQL Backup Task #" + smart.ConvertInt64ToStr(taskNum) + " :: " + DateTimeStartUtc + " #####"));
	fmt.Println("")
	//--

	//--
	if((smart.StrTrimWhitespaces(BKP_ARCHIVE_FOLDER) == "") ||
		(len(BKP_ARCHIVE_FOLDER) < 3) ||
		smart.StrContains(BKP_ARCHIVE_FOLDER, ".") ||
		!smart.StrRegexMatchString(smart.REGEX_SMART_SAFE_FILE_NAME, BKP_ARCHIVE_FOLDER)) { // dissalow path, just a folder name !
		//--
		logBackupError("Invalid Backup Archive Dir Name: `" + BKP_ARCHIVE_FOLDER + "`")
		uxmIsBackupRunning = false
		return
		//--
	} //end if
	//--
	isArchDirBkOk, errArchDirBk := reCreateDir(BKP_ARCHIVE_FOLDER)
	if((isArchDirBkOk != true) || (errArchDirBk != "")) {
		logBackupError("Backup Archive Folder Cleanup FAILED: # ", errArchDirBk)
		uxmIsBackupRunning = false
		return
	} //end if
	//--
	fmt.Println(color.YellowString("## Backup Archive Dir was cleared and re-created: `" + BKP_ARCHIVE_FOLDER + "` ##"))
	//--

	//--
	dtObjUtc := smart.DateTimeStructUtc("")
	var theArchName = BKP_ARCHIVE_FOLDER + "/" + "pg-db-dump-" + dtObjUtc.Years + dtObjUtc.Months + dtObjUtc.Days + "-" + dtObjUtc.Hours + dtObjUtc.Minutes + dtObjUtc.Seconds + ".7z"
	if((smart.StrTrimWhitespaces(theArchName) == "") ||
		smart.PathIsEmptyOrRoot(theArchName) ||
		smart.PathIsBackwardUnsafe(theArchName) ||
		smart.PathIsAbsolute(theArchName) ||
		smart.StrContains(theArchName, " ") ||
		!smart.StrRegexMatchString(smart.REGEX_SMART_SAFE_PATH_NAME, theArchName)) {
		//--
		logBackupError("Backup Archive File Path in Invalid or Unsafe: `" + theArchName + "`")
		uxmIsBackupRunning = false
		return
		//--
	} //end if
	//--

	//--
	var theBkpFolder = smart.PathDirName(BKP_SAFETY_FILE)
	if(smart.StrStartsWith(theBkpFolder, BKP_ARCHIVE_FOLDER)) {
		logBackupError("Backup Archive Dir `" + BKP_ARCHIVE_FOLDER + "` Name must be completely different than the Backup Folder Name `" + "`");
		uxmIsBackupRunning = false
		return
	} //end if
	//--
	isValid, errMsg := checkifSafeBackupDir(theBkpFolder, false)
	if((isValid != true) || (errMsg != "")) {
		logBackupError("Invalid Backup Dir: `" + theBkpFolder + "` # " + errMsg)
		uxmIsBackupRunning = false
		return
	} //end if
	if((theBkpFolder != smart.PathDirName(BKP_SCHEMA_FILE)) || (theBkpFolder != smart.PathDirName(BKP_DATA_FILE))) {
		logBackupError("Safety File, Schema File and DataFile must be all in the same Dir (currently detected: `" + theBkpFolder + "`)");
		uxmIsBackupRunning = false
		return
	} //end if
	//--
	isReDirBkOk, errReDirBk := reCreateDir(theBkpFolder)
	if((isReDirBkOk != true) || (errReDirBk != "")) {
		logBackupError("Backup Folder Cleanup FAILED: # ", errReDirBk)
		uxmIsBackupRunning = false
		return
	} //end if
	//--
	fmt.Println(color.YellowString("## Backup Dir was cleared and re-created: `" + theBkpFolder + "` ##"))
	//--

	//--
	areValid, errDet := checkIfSafePgDumpConnectionParams()
	if((areValid != true) || (errDet != "")) {
		logBackupError("Invalid DB Parameters:", errDet)
		uxmIsBackupRunning = false
		return
	} //end if
	//--

	//--
	isClean, errClean := cleanupOldFile(BKP_SAFETY_FILE)
	if((isClean != true) || (errClean != "")) {
		logBackupError("Cannot remove the Old Safety File (`" + BKP_SAFETY_FILE + "`): # " + errClean)
		uxmIsBackupRunning = false
		return
	} //end if
	//--

	//-- backup pgsql DB schema
	fmt.Println("")
	isBkpSchemaValid, errBkpSchemaMsg := doBackup("schema")
	if((isBkpSchemaValid != true) || (errBkpSchemaMsg != "")) {
		logBackupError("DB Schema Backup Failed:", errBkpSchemaMsg)
		uxmIsBackupRunning = false
		return
	} //end if
	//--

	//-- backup pgsql DB data
	fmt.Println("")
	isBkpDataValid, errBkpDataMsg := doBackup("data")
	if((isBkpDataValid != true) || (errBkpDataMsg != "")) {
		logBackupError("DB Data Backup Failed:", errBkpDataMsg)
		uxmIsBackupRunning = false
		return
	} //end if
	//--

	//--
	var DateTimeEndUtc string = smart.DateNowUtc()
	//--

	//--
	isSuccess, errMsg := smart.SafePathFileWrite("START @ " + DateTimeStartUtc + "\n" + "END   @ " + DateTimeEndUtc + "\n", "w", BKP_SAFETY_FILE, false)
	if((isSuccess != true) || (errMsg != "")) {
		logBackupError("Failed to create the after-backup Safety File: `" + BKP_SAFETY_FILE + "`", errMsg)
		uxmIsBackupRunning = false
		return
	} //end if
	//--

	//--
	fmt.Println("")
	fmt.Println(color.YellowString("## Backup Archive File is set to: `" + theArchName + "` ##"))
	fmt.Println("")
	//--
	fmt.Println("========== Archiving (7-Zip): START ==========")
	fmt.Println("7za: `" + theBkpFolder + "/` > `" + theArchName + "`")
	//--
	isSuccess, _, errStd := smart.ExecTimedCmd(CMD_TIMEOUT, "output", "capture+output", "", "", "7za", "a", "-t7z", "-m0=lzma", "-mx=5", "-md=256m", "-bb0", theArchName, theBkpFolder + "/")
	if((isSuccess != true) || (errStd != "")) {
		logBackupError("7-Zip Archiver encountered Errors / StdErr:\n`" + errStd + "`\n")
		uxmIsBackupRunning = false
		return
	} //end if
	if((!smart.PathExists(theArchName)) || (!smart.PathIsFile(theArchName))) {
		logBackupError("7-Zip Archive cannot be found (after backup + archiving) `" + theArchName + "`")
		uxmIsBackupRunning = false
		return
	} //end if
	fSize, errSize := createSzFile(theArchName)
	if((fSize <= 0) || (errSize != "")) {
		logBackupError("Failed to Create the Size Registration File of Backup Archive (after backup + archiving): `" + theArchName + "`, Size=" + smart.ConvertInt64ToStr(fSize) + " # " + errSize)
		uxmIsBackupRunning = false
		return
	} //end if
	//--
	fmt.Println("========== Archiving (7Zip): END ==========")
	//--

	//--
	fmt.Println("")
	//--
	isReClrDirBkOk, errReClrDirBk := smart.SafePathDirDelete(theBkpFolder, false)
	if((isReClrDirBkOk != true) || (errReClrDirBk != "") || (smart.PathExists(theBkpFolder))) {
		logBackupError("Backup Folder Cleanup (after archiving) FAILED: # ", errReClrDirBk)
		uxmIsBackupRunning = false
		return
	} //end if
	//--
	fmt.Println(color.YellowString("## Backup Dir was removed (after archiving, to save space): `" + theBkpFolder + "` ##"))
	//--

	//--
	fmt.Println("")
	fmt.Println(color.HiGreenString("### [OK: backup COMPLETED] :: " + DateTimeEndUtc + " ###"))
	fmt.Println("")
	//--

	//--
	uxmIsBackupRunning = false
	//--

} //END FUNCTION


//---


func main() {

	//--
	smart.LogToConsole("DEBUG", true) // log errors to console, with colors
	//--

	//--
	smart.ClearPrintTerminal()
	//--
	fmt.Println("")
	fmt.Println(color.HiBlueString("{# M-DataCenter :: PostgreSQL Backup (" + PROGR_VERSION + ") :: " + smart.DateNowUtc() + " :: (c) 2020 unix-world.org #}"));
	fmt.Println("")
	//--

	//--
//	go func() {
//		log.Fatal(http.Serve(ln, nil))
//	}()
	//--

	//-- main loop for task engine
	for i := 0; i >= 0; i++ { // infinite loop
		//--
		runBackupTask(int64(i+1))
		time.Sleep(time.Duration(30) * time.Second)
		//--
	}
	//--


} //END FUNCTION


//---


// #END
