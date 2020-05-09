
// GO Lang :: M-DataCenter :: Pg-Backup
// (c) 2020 unix-world.org
// r.20200509.1721 :: STABLE

package main


import (
	"os"
	"log"
	"fmt"

	smart "github.com/unix-world/smartgo"
)


const (
	CMD_TIMEOUT = 120

	PG_HOST = "127.0.0.1"
	PG_PORT = "5432"
	PG_USER = "pgsql"
	PG_PASS = "pgsql"
	PG_DB   = "smart_framework"

	BKP_FILE = "db-backup/t.sql"
)


func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION


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
	//--
	if((smart.StrTrimWhitespaces(fDirName) == "") ||
		smart.PathIsEmptyOrRoot(fDirName) ||
		smart.PathIsBackwardUnsafe(fDirName) ||
		smart.PathIsAbsolute(fDirName) ||
		!smart.PathExists(fDirName) ||
		!smart.PathIsDir(fDirName) ||
		!smart.StrRegexMatchString(smart.REGEX_SMART_SAFE_PATH_NAME, fDirName)) {
		//--
		return false, "Invalid or Non-Existing Backup Dir: " + fDirName
		//--
	} //end if
	//--
	var fBaseName string = smart.PathBaseName(bkpFile)
	//--
	if((smart.StrTrimWhitespaces(fBaseName) == "") ||
		!smart.StrEndsWith(fBaseName, ".sql") || (smart.PathBaseExtension(fBaseName) != ".sql") ||
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


func main() {

	//--
	smart.LogToConsole("DEBUG", true) // log errors to console, with colors
	//--

	//--
	areValid, errDet := checkIfSafePgDumpConnectionParams()
	if((areValid != true) || (errDet != "")) {
		fatalError("Invalid DB Parameters:", errDet)
	} //end if
	//--

	//--
	isValid, errMsg := checkIfSafeBackupFilePath(BKP_FILE)
	if((isValid != true) || (errMsg != "")) {
		fatalError("Invalid DB Backup File Path:", errMsg)
	} //end if
	//--

	//-- backup pgsql DB schema
	if(smart.PathExists(BKP_FILE)) {
		if(smart.PathIsDir(BKP_FILE)) {
			fatalError("Cannot do the backup. The backup file path is a Dir:", BKP_FILE)
		} //end if
		if(smart.PathIsFile(BKP_FILE)) {
			smart.SafePathFileDelete(BKP_FILE, false)
		} //end if
		if(smart.PathExists(BKP_FILE)) {
			fatalError("Cannot do the backup. The backup file path cannot be removed:", BKP_FILE)
		} //end if
	} //end if
	var pgDumpDetails string = "Host=" + PG_HOST + ":" + PG_PORT + " ; User=" + PG_USER + " ; Pass=***** ; DB=" + PG_DB + " ; File=" + BKP_FILE
	fmt.Println("========== PgDump SCHEMA: START ==========")
	fmt.Println("PgDump: SCHEMA # " + pgDumpDetails)
	isSuccess, _, errStd := smart.ExecTimedCmd(CMD_TIMEOUT, "output", "capture+output", "PGPASSWORD=" + PG_PASS, "", "pg_dump", "--encoding=UTF8", "--column-inserts", "--blobs", "--schema-only", "--no-owner" , "--no-privileges", "--host=" + PG_HOST, "--port=" + PG_PORT, "--user=" + PG_USER, "--format=p", "--file=" + BKP_FILE, PG_DB)
	if((isSuccess != true) || (errStd != "")) {
		fatalError("DB Backup Failed with Errors", "StdErr:\n`", errStd, "`\n")
	} //end if
	if((!smart.PathExists(BKP_FILE)) || (!smart.PathIsFile(BKP_FILE))) {
		fatalError("DB Backup Failed. Backup File cannot be found at:", BKP_FILE)
	} //end if
	fSize, errMsg := smart.SafePathFileGetSize(BKP_FILE, false)
	if(errMsg != "") {
		fatalError("DB Backup Failed. Cannot Get Size of the Backup File:", BKP_FILE, "#", errMsg)
	} //end if
	if(fSize <= 0) {
		fatalError("DB Backup Failed. Backup File is empty (", fSize, "bytes ):", BKP_FILE)
	} //end if
	//--
	fmt.Println("OK: DB SCHEMA was saved to:", BKP_FILE)
	fmt.Println("========== PgDump SCHEMA: END ==========")
	//--

} //END FUNCTION

// #END
