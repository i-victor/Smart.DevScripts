
// GO Lang
// sync / concurrency
// (c) 2017-2018 unix-world.org
// version: 2018.12.02

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {

	// concurency
//	runtime.GOMAXPROCS(1)

	// parallelism
	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Starting Go Routines")
	go func() {
		defer wg.Done()

		for char := 'a'; char < 'a'+26; char++ {
			time.Sleep(1 * time.Microsecond)
			fmt.Printf("%c ", char)
		}
	}()

	go func() {
		defer wg.Done()

		for number := 1; number < 27; number++ {
			fmt.Printf("%d ", number)
		}
	}()

	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")

}

// #END
