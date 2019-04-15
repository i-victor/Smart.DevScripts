
// GO Lang
// sync / semaphore
// (c) 2017-2019 unix-world.org
// version: 2019.04.15

package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.Ltime)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3)
	for i := 0; i < 9; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			semaphore <- struct{}{} // Lock
			defer func() {
				<-semaphore // Unlock
			}()

			// Resource intensive work goes here
			time.Sleep(time.Second)
			log.Println(i)
		}(i)
	}
	wg.Wait()
}

// #END
