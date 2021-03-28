
// GO Lang :: M-DataCenter :: Pg-Backup
// (c) 2020-2021 unix-world.org
// STABLE

package main


import (
	"os"
	"log"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"

	color "github.com/unix-world/smartgo/colorstring"
	smart "github.com/unix-world/smartgo"
)


const (
	PROGR_VERSION = "r.20210328.2258"

	CMD_TIMEOUT = 3600 // 1h per cmd
	CMD_RERUN_TIME = 30 // !IMPORTANT! This have to be greater than CMD_TIMEOUT

	PG_HOST = "127.0.0.1"
	PG_PORT = "5432"
	PG_USER = "pgsql"
	PG_PASS = "pgsql"
	PG_DB   = "smart_framework_r"

	PG_TEMPLATE_DB = "template0"
	PG_MASTER_DB = "postgres"

	BKP_RESTORE_FOLDER = "data-archive/" // this folder is where the 7zip archive must be found
	BKP_ARCHIVE_FOLDER = "data-backup/" // this folder must exists inside of the 7zip archive and here will be decompressed the archive
)


//--

type DbExists struct {
	Exists string `db:"the_db_exists"`
}

//--


func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION


//--


var uxmIsRestoreRunning bool = false


func logRestoreError(logMessages ...interface{}) {
	//--
	log.Println("[ERROR] :: RESTORE Task :: ", logMessages) // standard logger
	// must not exit
	//--
} //END FUNCTION


func logRestoreWarning(logMessages ...interface{}) {
	//--
	log.Println("[WARNING] :: RESTORE Task :: ", logMessages) // standard logger
	// must not exit
	//--
} //END FUNCTION


func logRestoreNotice(logMessages ...interface{}) {
	//--
	log.Println("[NOTICE] :: RESTORE Task :: ", logMessages) // standard logger
	// must not exit
	//--
} //END FUNCTION


//--


func getSzFileName(bkpFile string) string {
	//--
	var theSzFile string = bkpFile + ".size"
	//--
	return theSzFile
	//--
} //END FUNCTION


func checkFileSize(bkpFile string) (isOk bool, errMsg string) {
	//--
	var theSzFile string = getSzFileName(bkpFile)
	//--
	fSize, errMsg := smart.SafePathFileGetSize(bkpFile, false)
	if((errMsg != "") || (fSize <= 0)) {
		return false, "Could Not Get The Data Backup FileSize: `" + bkpFile + "` # " + errMsg
	} //end if
	//--
	fRealSize, errRealMsg := smart.SafePathFileRead(theSzFile, false)
	fRealSize = smart.StrTrimWhitespaces(fRealSize)
	if((errRealMsg != "") || (fRealSize == "")) {
		return false, "Could Not Read The Data Backup FileSize: `" + theSzFile + "` # " + errRealMsg
	} //end if
	//--
	var fStrSize string = smart.ConvertInt64ToStr(fSize)
	if(fStrSize != fRealSize) {
		return false, "Invalid Data Backup FileSize: " + fRealSize + " != " + fStrSize
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


//--


func getMd5FileName(bkpFile string) string {
	//--
	var theMd5ChskumFile string = bkpFile + ".md5"
	//--
	return theMd5ChskumFile
	//--
} //END FUNCTION


func checkFileMd5(bkpFile string) (isOk bool, errMsg string) {
	//--
	var theMd5File string = getMd5FileName(bkpFile)
	//--
	md5Sum, errMsg := smart.SafePathFileMd5(bkpFile, false)
	if((errMsg != "") || ((md5Sum == "") || (len(md5Sum) != 32))) {
		return false, "Failed to get the MD5 sum of File: `" + bkpFile + "` # " + errMsg
	} //end if
	//--
	fMd5ContentsRead, errMd5FileRead := smart.SafePathFileRead(theMd5File, false)
	fMd5ContentsRead = smart.StrTrimWhitespaces(fMd5ContentsRead)
	if((errMd5FileRead != "") || (fMd5ContentsRead == "")) {
		return false, "Could Not Read The Data Backup FileMd5: `" + theMd5File + "` # " + errMd5FileRead
	} //end if
	//--
	var theChecksumFileBaseName string = bkpFile
	if(smart.StrContains(theChecksumFileBaseName, "/")) {
		theChecksumFileBaseName = smart.PathBaseName(theChecksumFileBaseName)
	} //end if
	var expectedContent string = "MD5 (" + theChecksumFileBaseName + ") = " + md5Sum
	if(expectedContent != fMd5ContentsRead) {
		return false, "Invalid Data Backup FileMd5: `" + expectedContent + "` != `" + fMd5ContentsRead + "`"
	} //end if
	//--
	return true, ""
	//--
} //END FUNCTION


//--


func runRestoreTask(taskNum int64) {

	//--
	uxmIsRestoreRunning = true
	//--

	//--
	var DateTimeStartUtc string = smart.DateNowUtc()
	//--

	//--
	fmt.Println("")
	fmt.Println(color.BlueString("##### M-DataCenter :: PostgreSQL Restore Task #" + smart.ConvertInt64ToStr(taskNum) + " :: " + DateTimeStartUtc + " #####"));
	fmt.Println("")
	//--

	//--
	if(
		(smart.PathIsEmptyOrRoot(BKP_RESTORE_FOLDER)) ||
		(smart.PathIsAbsolute(BKP_RESTORE_FOLDER)) ||
		(smart.PathIsBackwardUnsafe(BKP_RESTORE_FOLDER)) ||
		(!smart.StrEndsWith(BKP_RESTORE_FOLDER, "/"))) {
			logRestoreError("Invalid Archive Dir: `" + BKP_RESTORE_FOLDER + "`")
			uxmIsRestoreRunning = false
			return
	} //end if
	//--
	if((!smart.PathExists(BKP_RESTORE_FOLDER)) || (!smart.PathIsDir(BKP_RESTORE_FOLDER))){
		// NOTICE, NOT ERROR !
		logRestoreNotice("Skip Task: The Archive Dir does NOT Exists")
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	if(
		(smart.PathIsEmptyOrRoot(BKP_ARCHIVE_FOLDER)) ||
		(smart.PathIsAbsolute(BKP_ARCHIVE_FOLDER)) ||
		(smart.PathIsBackwardUnsafe(BKP_ARCHIVE_FOLDER)) ||
		(!smart.StrEndsWith(BKP_ARCHIVE_FOLDER, "/"))) {
			logRestoreError("Invalid Backup Dir in Archive: `" + BKP_ARCHIVE_FOLDER + "`")
			uxmIsRestoreRunning = false
			return
	} //end if
	//--

	//--
	isReClrDirBkOk, errReClrDirBk := smart.SafePathDirDelete(BKP_ARCHIVE_FOLDER, false)
	if((isReClrDirBkOk != true) || (errReClrDirBk != "") || (smart.PathExists(BKP_ARCHIVE_FOLDER))) {
		logRestoreError("Restore Folder Cleanup (before unarchiving) FAILED: # ", errReClrDirBk)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--

	//--
	scanOk, errScanMsg, _, arrFiles := smart.SafePathDirScan(BKP_RESTORE_FOLDER, false, false)
	if((scanOk != true) || (errScanMsg != "")) {
		logRestoreError("Failed to Scan the Archive Dir: `" + BKP_RESTORE_FOLDER + "` :: " + errScanMsg)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--

	//--
	var archFileFound string = ""
	for _, file := range arrFiles {
		if(smart.StrEndsWith(file, ".7z")) {
			if(file > archFileFound) { // get the latest one
				archFileFound = file
			} //end if
		}
	} //end for
	//--
	if(
		(archFileFound == "") ||
		(!smart.PathExists(archFileFound)) ||
		(!smart.PathIsFile(archFileFound)) ||
		(!smart.StrEndsWith(archFileFound, ".7z"))) {
			logRestoreError("Failed to Find the 7Zip Archive File into: `" + BKP_RESTORE_FOLDER + "`")
			uxmIsRestoreRunning = false
			return
	} //end if
	archDoneFile := archFileFound + ".done-restore"
	if((smart.PathExists(archDoneFile)) || (smart.PathIsFile(archDoneFile))) {
		// NOTICE, NOT ERROR !
		logRestoreNotice("Skip Task: The Archive Was already processed: `" + archDoneFile + "`")
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	archSizeFileFound := archFileFound + ".size"
	//--
	if(
		(archSizeFileFound == "") ||
		(!smart.PathExists(archSizeFileFound)) ||
		(!smart.PathIsFile(archSizeFileFound)) ||
		(!smart.StrEndsWith(archSizeFileFound, ".7z.size"))) {
			logRestoreError("Failed to Find the 7Zip Archive Size File into: `" + BKP_RESTORE_FOLDER + "`")
			uxmIsRestoreRunning = false
			return
	} //end if
	//--

	//--
	fmt.Println(color.YellowString("## Backup Archive Dir was Found at: `" + BKP_RESTORE_FOLDER + "` ## :: " + archFileFound))
	//--

	//--
	theFSize, errFSizeMsg := smart.SafePathFileGetSize(archFileFound, false)
	if((theFSize <= 0) || (errFSizeMsg != "")) {
		logRestoreError("Failed to Get The FileSize of the 7Zip Archive File: `" + archFileFound + "` :: " + errFSizeMsg)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	theChkSize, errRdSizeFile := smart.SafePathFileRead(archSizeFileFound, false)
	theChkSize = smart.StrTrimWhitespaces(theChkSize)
	if(
		(theChkSize == "") ||
		(!smart.StrRegexMatchString(`^[0-9]+$`, theChkSize)) ||
		(errRdSizeFile != "")) {
			logRestoreError("Invalid 7Zip Archive Size File Content: `" + archSizeFileFound + "` :: " + errRdSizeFile)
			uxmIsRestoreRunning = false
			return
	} //end if
	//--
	var theIChkSize int64 = smart.ParseIntegerStrAsInt64(theChkSize)
	if(theIChkSize <= 0) {
		logRestoreError("Invalid 7Zip Archive Size File Int64 Content: `" + archSizeFileFound + "` ::", theIChkSize)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	if(theIChkSize != theFSize) {
		logRestoreError("Invalid 7Zip Archive File Size Check #", theFSize, "::", theIChkSize)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--

	//--
	fmt.Println("")
	fmt.Println(color.HiGreenString("### [OK: Archive Size Check] :: " + theChkSize + " Bytes ###"))
	fmt.Println("")
	//--

	//--
	fmt.Println("========== Un-Archiving (7-Zip): START ==========")
	fmt.Println("7za e: `" + archFileFound + "` > `" + BKP_RESTORE_FOLDER + "`")
	//--
	isSuccess, _, errStd := smart.ExecTimedCmd(CMD_TIMEOUT, "output", "capture+output", "", "", "7za", "e", "-y", archFileFound, "-o" + BKP_ARCHIVE_FOLDER, BKP_ARCHIVE_FOLDER + "*", "-r", "-bb1")
	if((isSuccess != true) || (errStd != "")) {
		logRestoreError("7-Zip Archiver encountered Errors / StdErr:\n`" + errStd + "`\n")
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	fmt.Println("========== Un-Archiving (7Zip): END ==========")
	//--
	if((!smart.PathExists(BKP_ARCHIVE_FOLDER)) || (!smart.PathIsDir(BKP_RESTORE_FOLDER))) {
		logRestoreError("Unarchive FAILED: The Data Backup Folder cannot be found after unarchiving: `" + BKP_ARCHIVE_FOLDER + "`")
		uxmIsRestoreRunning = false
		return
	} //end if
	//--

	//--
	fmt.Println("")
	//--

	//--
	fmt.Println("========== Checking Files Integrity: START ==========")
	//--
	if(
		(!smart.PathIsFile(BKP_ARCHIVE_FOLDER + "db-schema.sql")) ||
		(!smart.PathIsFile(BKP_ARCHIVE_FOLDER + "db-schema.sql.md5")) ||
		(!smart.PathIsFile(BKP_ARCHIVE_FOLDER + "db-schema.sql.size")) ||
		(!smart.PathIsFile(BKP_ARCHIVE_FOLDER + "db-data.sql")) ||
		(!smart.PathIsFile(BKP_ARCHIVE_FOLDER + "db-data.sql.md5")) ||
		(!smart.PathIsFile(BKP_ARCHIVE_FOLDER + "db-data.sql.size")) ||
		(!smart.PathIsFile(BKP_ARCHIVE_FOLDER + "db-pgdump.ok"))) {
			logRestoreError("Unarchive FAILED: The Data Backup Folder is missing some required files after unarchiving: `" + BKP_ARCHIVE_FOLDER + "`")
			uxmIsRestoreRunning = false
			return
	} //end if
	//--
	chkSchemaFSize, errSchemaFSize := checkFileSize(BKP_ARCHIVE_FOLDER + "db-schema.sql")
	if(chkSchemaFSize != true) {
		logRestoreError("Unarchive FAILED (Schema): " + errSchemaFSize)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	chkDataFSize, errDataFSize := checkFileSize(BKP_ARCHIVE_FOLDER + "db-data.sql")
	if(chkDataFSize != true) {
		logRestoreError("Unarchive FAILED (Data): " + errDataFSize)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	chkSchemaFMd5, errSchemaFMd5 := checkFileMd5(BKP_ARCHIVE_FOLDER + "db-schema.sql")
	if(chkSchemaFMd5 != true) {
		logRestoreError("Unarchive FAILED (Schema): " + errSchemaFMd5)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	chkDataFMd5, errDataFMd5 := checkFileMd5(BKP_ARCHIVE_FOLDER + "db-data.sql")
	if(chkDataFMd5 != true) {
		logRestoreError("Unarchive FAILED (Data): " + errDataFMd5)
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	fmt.Println(color.GreenString("### OK: Control File, Sizes and MD5 Checksums"))
	//--
	fmt.Println("========== Checking Files Integrity: END ==========")
	fmt.Println("")
	//--

	//--
	var maxConnTimeOut int = CMD_RERUN_TIME - 5
	if(maxConnTimeOut < 5) {
		maxConnTimeOut = 5
	} //end if
	connectTimeOut := time.Duration(maxConnTimeOut) * time.Second
	//--
	dtObjUtc := smart.DateTimeStructUtc("")
	var theDbName = PG_DB + "_" + dtObjUtc.Years + "_" + dtObjUtc.Months + "_" + dtObjUtc.Days
	fmt.Println("========== DB Restore: START [" + theDbName + "] ==========")
	fmt.Println("")
	//--
	var theCreateDbSQL string = "CREATE DATABASE " + theDbName + " OWNER " + PG_USER + " ENCODING 'UTF8' LC_COLLATE 'C' TEMPLATE " + PG_TEMPLATE_DB + ";"
	var theSetDbOwnerSQL string = "ALTER SCHEMA public OWNER TO " + PG_USER + ";"
	//--
	dbMasterConn, errMasterConn := sqlx.Connect("postgres", "host=" + PG_HOST + " port=" + PG_PORT + " user=" + PG_USER + " password=" + PG_PASS + " dbname=" + PG_MASTER_DB + " sslmode=disable")
	if(errMasterConn != nil) {
		logRestoreError("Failed to Connect to PostgreSQL Server / Master DB: " + errMasterConn.Error())
		uxmIsRestoreRunning = false
		return
	} //end if
	dbMasterConn.SetConnMaxLifetime(connectTimeOut)
	//--
	resultDbExists := []DbExists{}
	err := dbMasterConn.Select(&resultDbExists, "SELECT 1 AS the_db_exists FROM pg_database WHERE datname = $1 LIMIT 1 OFFSET 0", theDbName)
	if(err != nil) {
		log.Fatalln(err)
		return
	}
	if(len(resultDbExists) > 0) {
		if(resultDbExists[0].Exists == "1") {
			logRestoreError("The PostgreSQL Server DB `" + theDbName + "` already exists ...")
			uxmIsRestoreRunning = false
			return
		} //end if
	} //end if
	//--
	_, errCreateDb := dbMasterConn.Exec(theCreateDbSQL)
	if(errCreateDb != nil) {
		logRestoreError("DB Create ERRORS: " + errCreateDb.Error())
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	dbConn, errDbConn := sqlx.Connect("postgres", "host=" + PG_HOST + " port=" + PG_PORT + " user=" + PG_USER + " password=" + PG_PASS + " dbname=" + theDbName + " sslmode=disable")
	if(errDbConn != nil) {
		logRestoreError("Failed to Connect to PostgreSQL Server / DB [" + theDbName + "]: " + errDbConn.Error())
		uxmIsRestoreRunning = false
		return
	} //end if
	dbConn.SetConnMaxLifetime(connectTimeOut)
	//--
	_, errSetPrivDb := dbConn.Exec(theSetDbOwnerSQL)
	if(errSetPrivDb != nil) {
		logRestoreError("DB Set Privileges ERRORS: " + errSetPrivDb.Error())
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	isRestoreSchemaSuccess, _, errRestoreSchemaStd := smart.ExecTimedCmd(CMD_TIMEOUT, "output", "capture+output", "PGPASSWORD=" + PG_PASS, "", "psql", "--quiet", "--echo-errors", "--host=" + PG_HOST, "--port=" + PG_PORT, "--username=" + PG_USER, "--no-password", "--dbname=" + theDbName, "--file=" + BKP_ARCHIVE_FOLDER + "db-schema.sql")
	if((isRestoreSchemaSuccess != true) || (errRestoreSchemaStd != "")) {
		logRestoreError("DB Restore FAILED (Schema): Errors / StdErr:\n`" + errRestoreSchemaStd + "`\n")
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	isRestoreDataSuccess, _, errRestoreDataStd := smart.ExecTimedCmd(CMD_TIMEOUT, "output", "capture+output", "PGPASSWORD=" + PG_PASS, "", "psql", "--quiet", "--echo-errors", "--host=" + PG_HOST, "--port=" + PG_PORT, "--username=" + PG_USER, "--no-password", "--dbname=" + theDbName, "--file=" + BKP_ARCHIVE_FOLDER + "db-data.sql")
	if((isRestoreDataSuccess != true) || (errRestoreDataStd != "")) {
		logRestoreError("DB Restore FAILED (Data): Errors / StdErr:\n`" + errRestoreDataStd + "`\n")
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	fmt.Println("========== DB Restore: END [" + theDbName + "] ==========")
	fmt.Println("")
	//--

	//--
	doneFileWrOk, doneFileWrErr := smart.SafePathFileWrite("DONE: " + smart.DateNowUtc(), "w", archDoneFile, false)
	if((doneFileWrOk != true) || (doneFileWrErr != "")) {
		logRestoreWarning("FAILED to create the Control File (" + archDoneFile + ") as: `" + doneFileWrErr + "`\n")
		uxmIsRestoreRunning = false
		return
	} //end if
	//--
	fmt.Println(color.YellowString("## DONE: Control File Created as: `" + archDoneFile + "` ##"))
	//--

	//--
	var DateTimeEndUtc string = smart.DateNowUtc()
	//--

	//--
	fmt.Println("")
	fmt.Println(color.HiGreenString("### [OK: Restore COMPLETED using PostgreSQL Database `" + theDbName + "`] :: " + DateTimeEndUtc + " ###"))
	fmt.Println("")
	//--

	//--
	uxmIsRestoreRunning = false
	//--

} //END FUNCTION


func main() {

	//--
	smart.LogToConsole("DEBUG", true) // log errors to console, with colors
	//--

	//--
	smart.ClearPrintTerminal()
	//--
	fmt.Println("")
	fmt.Println(color.HiBlueString("{# M-DataCenter :: PostgreSQL Restore (" + PROGR_VERSION + ") :: " + smart.DateNowUtc() + " :: (c) 2020 unix-world.org #}"));
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
		runRestoreTask(int64(i+1))
		time.Sleep(time.Duration(CMD_RERUN_TIME) * time.Second)
		//--
	}
	//--

} //END FUNCTION


//---


// #END

