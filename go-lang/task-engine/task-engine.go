
// GO Lang

// parallel task engine
// (c) 2018 unix-world.org

package main

import (
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"time"
	"math/rand"
	"io/ioutil"
	"net/http"
	"net/url"
	"crypto/tls"
	"runtime"
	"sync"
	"os"
	"os/signal"
	"syscall"
//	"reflect"

	"github.com/vaughan0/go-ini"
	"github.com/fatih/color"
)

var uxmScriptVersion = "r.180227"

var iniFile = "task-engine.ini"
var UrlBatchList = "https://localhost/test-batch.txt"
var UrlTaskCall = "https://localhost/task?TaskID="

var parallelWorkers = 2

func tasks(maxParallelThreads int, LoopId int) {

	fmt.Println(color.HiMagentaString("@@@@@ Max Parallel Threads: " + strconv.Itoa(maxParallelThreads) + " @@@@@"))
	time.Sleep(time.Duration(1) * time.Second)

	fmt.Println("===== Get Batch =====")

	if(UrlBatchList == "") {
		fmt.Println(color.RedString("ERROR: Batch List URL is Empty"))
		return
	}
	if(UrlTaskCall == "") {
		fmt.Println(color.RedString("ERROR: Task Call URL is Empty"))
		return
	}

	fmt.Println(color.MagentaString("Getting the Batch List from URL: " + UrlBatchList))
	res, err := http.Get(UrlBatchList)
	if err != nil {
		fmt.Println(color.RedString("ERROR accessing Batch List URL"))
		fmt.Println(err.Error())
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(color.RedString("ERROR reading Batch List Body"))
		fmt.Println(err.Error())
		return
	}
	if(res.StatusCode != 200) {
		fmt.Println(color.RedString("ERROR reading Batch List :: Status Code = " + strconv.Itoa(res.StatusCode)))
		return
	}
	list := strings.Trim(string(body), " \t\n\r\x00\x0B");
	if(list == "") {
		fmt.Println(color.RedString("ERROR processing Batch List :: List is Empty"))
		return
	}
	arr := strings.Split(list, "\n")
	if(len(arr) <= 0) {
		fmt.Println(color.RedString("ERROR processing Batch List :: No List Entries Found"))
		return
	}
	if(arr[0] != "#IDs-BATCH:START#") {
		fmt.Println(color.RedString("ERROR processing Batch List :: Invalid Starting Line"))
		return
	}
	if(arr[len(arr)-1] != "#IDs-BATCH:END#") {
		fmt.Println(color.RedString("ERROR processing Batch List :: Invalid Ending Line"))
		return
	}

	fmt.Println(color.MagentaString("Task Processing URL: " + UrlTaskCall))

	var wg sync.WaitGroup

	for i := range arr {
		var TaskId = strings.Trim(arr[i], " \t\n\r\x00\x0B")
		var validID = regexp.MustCompile(`^[_A-Za-z0-9\-]*$`)
		if(TaskId == "#IDs-BATCH:START#") {
			fmt.Println(color.BlueString("***** Skip Pre-Processing: Batch#Start *****"))
			continue // skip
		}
		if(TaskId == "#IDs-BATCH:END#") {
			fmt.Println(color.BlueString("***** Skip Pre-Processing: Batch#End *****"))
			continue // skip
		}
		if(TaskId == "") {
			fmt.Println(color.YellowString("***** Skip Pre-Processing: Empty Task ID *****"))
			continue // skip
		}
		if(!validID.MatchString(TaskId)) {
			fmt.Println(color.YellowString("***** Skip Pre-Processing: Invalid Task ID: " + TaskId + " *****"))
			continue // skip
		}
		wg.Add(1)
		go func(LoopId int, i int, TaskId string) {
			res, err := http.Get(UrlTaskCall + url.QueryEscape(TaskId))
			if err != nil {
				fmt.Println(err.Error())
			}
			res.Body.Close()
			var status = ""
			if(res.StatusCode == 200) {
				status = color.GreenString(strconv.Itoa(res.StatusCode))
			} else if(res.StatusCode == 202) {
				status = color.HiCyanString(strconv.Itoa(res.StatusCode))
			} else if(res.StatusCode == 203) {
				status = color.YellowString(strconv.Itoa(res.StatusCode))
			} else if(res.StatusCode == 208) {
				status = color.HiMagentaString(strconv.Itoa(res.StatusCode))
			} else if(res.StatusCode == 429) || (res.StatusCode == 502) || (res.StatusCode == 503) || (res.StatusCode == 504) {
				status = color.HiRedString(strconv.Itoa(res.StatusCode))
			} else {
				status = color.RedString(strconv.Itoa(res.StatusCode))
			}
			fmt.Println("Task # " + color.HiYellowString(TaskId) + color.CyanString(" @ Thread.ID:")  + color.HiBlackString(strconv.Itoa(LoopId)) + "." + color.HiBlueString(strconv.Itoa(i)) + " :: HTTP Response Status:" + status)
			wg.Done()
		}(LoopId, i, TaskId)
		// just add a random pause in milliseconds for give a small breath ... (important for spread of threads in time !!!)
		time.Sleep(time.Duration(rand.Int31n(3) * 25) * time.Millisecond)
	}

	wg.Wait()

	fmt.Println("===== Done Processing Batch Tasks =====")

}

func getIniStrVal(file ini.File, Section string, Key string) string {
	str, ok := file.Get(Section, Key)
	if !ok {
		str = ""
	}
	return str
}

func getIniIntVal(file ini.File, Section string, Key string) int {
	str := getIniStrVal(file, Section, Key)
	i, err := strconv.Atoi(str)
	if err != nil {
		i = 0
	}
	return i
}

func main() {

	file, err := ini.LoadFile(iniFile)
//	fmt.Printf(reflect.TypeOf(file).String())
	if err != nil {
		fmt.Println(color.RedString("ERROR :: INI File cannot be loaded: " + iniFile))
		return
	}

	UrlBatchList = getIniStrVal(file, "URLs", "batch-url")
	if(UrlBatchList == "") {
		fmt.Println(color.RedString("ERROR :: Invalid key [URLs/batch-url] in INI File: " + iniFile))
		return
	}
	UrlTaskCall = getIniStrVal(file, "URLs", "tasks-url")
	if(UrlTaskCall == "") {
		fmt.Println(color.RedString("ERROR :: Invalid key [URLs/tasks-url] in INI File: " + iniFile))
		return
	}
	parallelWorkers = getIniIntVal(file, "Tunnings", "parallel-workers")
	if (parallelWorkers < 2) || (parallelWorkers > 16384) {
		parallelWorkers = runtime.NumCPU()
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(color.HiBlackString("\n" + "##### DONE Parallel Tasks Manager #####"))
		os.Exit(1)
	}()

	fmt.Println(color.HiBlackString("\n" + "##### Parallel Tasks Manager [ " + uxmScriptVersion + " ] :: Running on #CPUs: " + strconv.Itoa(runtime.NumCPU()) + " #####"))

	var maxParallelThreads = parallelWorkers

	runtime.GOMAXPROCS(maxParallelThreads)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	for i := 0; i >= 0; i++ { // infinite loop
		tasks(maxParallelThreads, i)
		time.Sleep(time.Duration(rand.Int31n(2)) * time.Second)
	}

}

//#END
