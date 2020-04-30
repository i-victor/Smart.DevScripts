
// GO Lang
// progress bar
// (c) 2020 unix-world.org

package main

import (
	"fmt"
	"time"
	"github.com/schollz/progressbar"
)

func main() {

	//bar := progressbar.NewOptions(100, progressbar.OptionSetRenderBlankState(true))
	bar := progressbar.New(100)
	bar.RenderBlank()

	for i := 0; i < 100; i++ {
		bar.Set(i)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n" + "The End!")

}

// #END
