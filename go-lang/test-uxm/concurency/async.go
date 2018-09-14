
package main

import (
	"fmt"
	"github.com/ti/nasync"
)

func main() {
		//new a async pool in max 1000 task in max 1000 gorutines
		async := nasync.New(1000,1000)
		defer async.Close()
		async.Do(doSometing,"hello word")
}


func doSometing(msg string) string{
	fmt.Println("i am done by " + msg)
	return "i am done by " + msg
}

