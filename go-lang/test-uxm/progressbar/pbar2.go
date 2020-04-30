
// GO Lang
// progress bar, alternative
// (c) 2020 unix-world.org

package main

import (
	"time"
	"github.com/cheggaaa/pb"
)

func main() {

	count := 500
	bar := pb.New(count)

	// show percents (by default already true)
	bar.ShowPercent = true

	// show bar (by default already true)
	bar.ShowBar = true

	bar.ShowCounters = true

	bar.ShowTimeLeft = true

	// and start
	bar.Start()
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("The End!")

}

// #END
