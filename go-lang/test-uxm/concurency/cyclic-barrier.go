
// GO Lang
// cyclic barrier
// (c) 2020 unix-world.org
// version: 2020.02.21

package main

import (
	"fmt"
	"sync"
	"time"
	"context"
	"github.com/marusama/cyclicbarrier"
)

func main() {

    // create a barrier for 10 parties with an action that increments counter
    // this action will be called each time when all goroutines reach the barrier
    cnt := 0
    b := cyclicbarrier.NewWithAction(10, func() error {
	cnt++
	return nil
    })

    wg := sync.WaitGroup{}
    for i := 0; i < 10; i++ { // create 10 goroutines (the same count as barrier parties)
	wg.Add(1)
	go func() {
	    for j := 0; j < 5; j++ {
		// do some hard work 5 times
		time.Sleep(100 * time.Millisecond)
		err := b.Await(context.TODO()) // ..and wait for other parties on the barrier. Last arrived goroutine will do the barrier action and then pass all other goroutines to the next round
		if err != nil {
		    panic(err)
		}
	    }
	    wg.Done()
	}()
    }

    wg.Wait()
    fmt.Println(cnt) // cnt = 5, it means that the barrier was passed 5 times

}

// #END
