
// GO Lang
// sync parallel / concurrency
// (c) 2017-2021 unix-world.org
// version: 20210328.2258

package main

import (
	"fmt"
	"strconv"
	"net/http"
	"crypto/tls"
	"runtime"
	"sync"
	"time"
	"math/rand"

	"github.com/fatih/color"
//	color "github.com/unix-world/smartgo/colorstring"
)

func main() {

	start := time.Now()

	fmt.Println("Running on #CPUs: " + strconv.Itoa(runtime.NumCPU()))

	runtime.GOMAXPROCS(4)

	var wg sync.WaitGroup

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var c = 0; // for c := 0; c < 2; c++ {
	for i := 0; i < 800; i++ {
		wg.Add(1)
		go func() {
			res, err := http.Get("https://127.0.0.1/sites/")
			if err != nil {
				fmt.Println(color.RedString("ERROR accessing URL"))
			//	fmt.Println(err)
				fmt.Println(err.Error())
			} else {
				res.Body.Close()
				var status = ""
				if(res.StatusCode == 200) {
					status = color.GreenString(strconv.Itoa(res.StatusCode))
				} else {
					status = color.RedString(strconv.Itoa(res.StatusCode))
				}
				fmt.Println("HTTP Response [" + strconv.Itoa(c) + "/"  + strconv.Itoa(i) +  "] Status: " + status)
			//	time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
			}
			wg.Done()
		}()
		time.Sleep(time.Duration(rand.Int31n(2) * 25) * time.Millisecond)
	}
//	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println("Time elapsed: %s", elapsed)

}

// #END
