
// Go Lang
// Date and Time samples
// r.20200505.2315 :: STABLE

package main

import (
	"fmt"

	smart "github.com/unix-world/smartgo"
)

func main() {


	fmt.Println("DateTime ISO + TZ-OFS (now) UTC is: ",  smart.DateNowUtc())
	fmt.Println("DateTime ISO (now) UTC is: ",  smart.DateNowIsoUtc())
	fmt.Println("DateTime ISO + TZ-OFS (now) LOCAL is: ", smart.DateNowLocal())
	fmt.Println("DateTime ISO (now) LOCAL is: ", smart.DateNowIsoLocal())

//	dTimeStr := "2020-01-05 08:03:07 +0300"
	dTimeStr := ""

	dtObjUtc := smart.DateTimeStructUtc(dTimeStr)
	fmt.Println("Date Obj Json", smart.JsonEncode(dtObjUtc))

	dtObjLoc := smart.DateTimeStructLocal(dTimeStr)
	fmt.Println("Date Obj Json", smart.JsonEncode(dtObjLoc))

//	fmt.Println("Year (now) UTC is: ", dtObj.Year)
//	fmt.Println("Month (now) UTC is: ", dtObj.Month)

}

// #END
