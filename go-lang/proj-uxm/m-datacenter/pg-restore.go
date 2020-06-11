
// GO Lang :: M-DataCenter :: Pg-Backup
// (c) 2020 unix-world.org
// STABLE

package main


import (
	"os"
	"log"
	"fmt"

	"github.com/fatih/color"
	smart "github.com/unix-world/smartgo"
)


const (
	PROGR_VERSION = "r.20200510.1527"

	CMD_TIMEOUT = 120
)


//---


func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	os.Exit(1)
	//--
} //END FUNCTION


//--


func main() {

	//--
	smart.LogToConsole("DEBUG", true) // log errors to console, with colors
	//--

	//--
	smart.ClearPrintTerminal()
	//--

	//--
	// PGPASSWORD=pgsql psql --quiet --echo-errors --host=127.0.0.1 --port=5432 --username=pgsql --no-password --dbname=m_datacenter_test --file=db-schema.sql
	//--

	//--
	var DateTimeEndUtc string = smart.DateNowUtc()
	//--

	//--
	fmt.Println("")
//	fmt.Println(color.HiGreenString("### [OK: restore COMPLETED] :: " + DateTimeEndUtc + " ###"))
	fmt.Println(color.RedString("### [ERR: restore NOT YET IMPLEMENTED] :: " + DateTimeEndUtc + " ###"))
	fmt.Println("")
	//--


} //END FUNCTION


//---


// #END

