
// GO Lang

// parallel task engine with built-in HTTP service control
// (c) 2017-2020 unix-world.org

package main


import (
	"os"
	"os/signal"
	"syscall"
	"sync"
	"runtime"
	"log"
	"fmt"
	"flag"
	"regexp"
	"strings"
	"strconv"
	"bytes"
	"time"
	"math/rand"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"html"

	"github.com/vaughan0/go-ini"
	"github.com/fatih/color"
)


const (
	THE_VERSION = "r.20200425.1635"
	INI_FILE = "task-engine.ini"
	SVG_LOGO = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16"><path fill-rule="evenodd" d="M3 4.949a2.5 2.5 0 10-1 0v8.049c0 .547.453 1 1 1h2.05a2.5 2.5 0 004.9 0h1.1a2.5 2.5 0 100-1h-1.1a2.5 2.5 0 00-4.9 0H3v-5h2.05a2.5 2.5 0 004.9 0h1.1a2.5 2.5 0 100-1h-1.1a2.5 2.5 0 00-4.9 0H3v-2.05zm9 2.55a1.5 1.5 0 103 0 1.5 1.5 0 00-3 0zm-3 0a1.5 1.5 0 10-3 0 1.5 1.5 0 003 0zm4.5 7.499a1.5 1.5 0 110-3.001 1.5 1.5 0 010 3zm-6-3a1.5 1.5 0 110 3 1.5 1.5 0 010-3z"/></svg>`
	SVG_SPIN = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32" fill="grey" id="loading-spin-svg"><path opacity=".25" d="M16 0 A16 16 0 0 0 16 32 A16 16 0 0 0 16 0 M16 4 A12 12 0 0 1 16 28 A12 12 0 0 1 16 4"/><path d="M16 0 A16 16 0 0 1 32 16 L28 16 A12 12 0 0 0 16 4z"><animateTransform attributeName="transform" type="rotate" from="0 16 16" to="360 16 16" dur="0.8s" repeatCount="indefinite" /></path></svg>`
	HTML_START = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>TaskEngine GO MiniServer</title>
<style>
pre#json-code { font-weight:bold; font-size:1rem; }
pre#json-code .string { color: #003399; }
pre#json-code .number { color: #FF3300; }
pre#json-code .boolean { color: #00CC00; }
pre#json-code .null { color: #9900FF; }
pre#json-code .key { color: #333333; }
</style>
</head>
<body>
`
	HTML_END = `
</body>
</html>
`
	HTML_JS = `
<script>
function jsonStrUnquote(str) {
	return str.replace(/("\:$)/g, ':').replace(/(^")|("$)/g, '').replace(/(\\")/g, '"');
}
function syntaxHighlight(json) {
	json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
	return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function(match) {
		var cls = 'number';
		if(/^"/.test(match)) {
			if(/:$/.test(match)) {
				cls = 'key';
			} else {
				cls = 'string';
			}
		} else if(/true|false/.test(match)) {
			cls = 'boolean';
		} else if(/null/.test(match)) {
			cls = 'null';
		}
		return '<span class="' + cls + '">' + jsonStrUnquote(match).replace(/"/g, '&quot;') + '</span>';
	});
}
function displayJson(elemID) {
	var str;
	var elObj;
	try {
		elObj = document.getElementById(elemID);
	} catch(err){
		return;
	}
	if(!elObj) {
		return;
	}
	try {
		str = elObj.innerText;
	} catch(err){
		return;
	}
	try {
		str = JSON.parse(str);
	} catch(err){
		return;
	}
	str = JSON.stringify(str, null, 4);
	elObj.innerHTML = syntaxHighlight(str);
}
try{ displayJson('json-code'); } catch(err){ console.log('FAILED to Color Format Json: ', err); }
</script>
`
)


var UrlBatchList string = ""
var UrlTaskCall string = ""
var bindTcpAddr string = ""

var parallelWorkers int = 2
var srvHostAddr string = "127.0.0.1"
var srvHostPort int = 0

var startTime string = ""

// in structures the keys that start with lowercase are private if used for json export ; they need to be re-mapped to json keys if need to have lowerkeys in json
type uxmStuctStats struct {
	Description string    `json:"description"`
	ServiceVersion string `json:"serviceVersion"`
	IsActive bool         `json:"isActive"`
	IsSilent bool         `json:"isSilent"`
	StartTime string      `json:"startTime"`
	CurrentTime string    `json:"currentTime"`
	BatchURL string       `json:"batchURL"`
	TaskURL string        `json:"taskURL"`
	NumWorkers int        `json:"numWorkers"`
	EmptyBatches int      `json:"emptyBatches"`
	BatchCycles int       `json:"batchCycles"`
	TasksOnWait int       `json:"tasksOnWait"`
	TasksProcessed int    `json:"tasksProcessed"`
	Stat200 uint64        `json:"stat200"`
	Stat202 uint64        `json:"stat202"`
	Stat203 uint64        `json:"stat203"`
	Stat208 uint64        `json:"stat208"`
	Stat429 uint64        `json:"stat429"`
	Stat502 uint64        `json:"stat502"`
	Stat503 uint64        `json:"stat503"`
	Stat504 uint64        `json:"stat504"`
	StatERR uint64        `json:"statERR"`
	StatALL uint64        `json:"statALL"`
}
var uxmStats *uxmStuctStats
var flagResetStats bool = false
var flagSilent bool = false


func tasks(maxParallelThreads int, LoopId int) {

	//--
	fmt.Println(color.HiMagentaString("----- Max Parallel Threads Execution: " + strconv.Itoa(maxParallelThreads) + " # ASYNC # -----"))
	fmt.Println(color.HiYellowString("INFO: Built-in MiniServer is listening at http://" + bindTcpAddr + "/"))
	//--

	//-- reset statistics if flagResetStats is set to true ...
	if(flagResetStats == true) {
		//--
		flagResetStats = false
		//--
		resetStatistics()
		//--
	}
	//--

	//--
	if(uxmStats.IsActive != true) {
		fmt.Println(color.HiYellowString("INFO: Service is PAUSED by Remote Control via HTTP"))
		if(flagSilent != true) {
			fmt.Println(color.GreenString("To Re-Activate the Service make a hit the following URL: " + "http://" + srvHostAddr + ":" + strconv.Itoa(srvHostPort) + "/start"))
		}
		time.Sleep(time.Duration(10) * time.Second)
		return
	} else {
		time.Sleep(time.Duration(1) * time.Second)
	}
	//--

	//--
	fmt.Println("===== Get Batch =====")
	uxmStats.BatchCycles++
	//--
	if(UrlBatchList == "") {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR: Batch List URL is Empty"))
		return
	}
	if(UrlTaskCall == "") {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR: Task Call URL is Empty"))
		return
	}
	//--

	//--
	fmt.Println(color.MagentaString("Getting the Batch List from URL: " + UrlBatchList))
	res, err := http.Get(UrlBatchList)
	if err != nil {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR accessing Batch List URL"))
		fmt.Println(err.Error())
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR reading Batch List Body"))
		fmt.Println(err.Error())
		return
	}
	if(res.StatusCode != 200) {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR reading Batch List :: Status Code = " + strconv.Itoa(res.StatusCode)))
		return
	}
	list := strTrimWhitespaces(string(body))
	if(list == "") {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR processing Batch List :: List is Empty"))
		return
	}
	arr := strings.Split(list, "\n")
	if(len(arr) <= 0) {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR processing Batch List :: No List Entries Found"))
		return
	}
	if(arr[0] != "#IDs-BATCH:START#") {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR processing Batch List :: Invalid Starting Line"))
		return
	}
	if(arr[len(arr)-1] != "#IDs-BATCH:END#") {
		uxmStats.StatERR++
		uxmStats.EmptyBatches++
		fmt.Println(color.RedString("ERROR processing Batch List :: Invalid Ending Line"))
		return
	}
	//--
	if(flagSilent != true) {
		fmt.Println(color.MagentaString("Task Processing URL: " + UrlTaskCall))
	}
	//--

	//--
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxParallelThreads)
	//--

	//-- strategy: will lauch twice the number of allowed max parallel tasks at once and later wait for pool by comparing max allowed with tasks on wait
	var maxTasksToLaunchLimit int = maxParallelThreads * 2 // {{{SYNC-TASKENGINE-FORK-LIMIT-PROTECTION}}} avoid launch more than 2 x maxParallelThreads to avoid system crash
	//--
	if(flagSilent != true) {
		fmt.Println(color.HiMagentaString("Parallel Threads Launch Limit: " + strconv.Itoa(maxTasksToLaunchLimit) + " # SYNC"))
	}
	//--

	//--
	var inListTasks int = 0
	//--
	for i := range arr {
		var TaskId = strTrimWhitespaces(string(arr[i]))
		if(TaskId == "#IDs-BATCH:START#") {
			if(flagSilent != true) {
				fmt.Println(color.BlueString("***** Skip Pre-Processing: Batch#Start *****"))
			}
			continue // skip
		}
		if(TaskId == "#IDs-BATCH:END#") {
			if(flagSilent != true) {
				fmt.Println(color.BlueString("***** Skip Pre-Processing: Batch#End *****"))
			}
			continue // skip
		}
		if(TaskId == "") {
			if(flagSilent != true) {
				fmt.Println(color.YellowString("***** Skip Pre-Processing: Empty Task ID *****"))
			}
			continue // skip
		}
		var validID = regexp.MustCompile(`^[_a-zA-Z0-9\-\.@]+$`) // allow all safe names except #
		if(!validID.MatchString(TaskId)) {
			uxmStats.StatALL++
			uxmStats.StatERR++
			fmt.Println(color.YellowString("***** Skip Pre-Processing: Invalid Task ID: " + TaskId + " *****"))
			continue // skip
		}
		inListTasks++
		//--
		if(uxmStats.TasksOnWait >= maxTasksToLaunchLimit) { // {{{SYNC-TASKENGINE-FORK-LIMIT-PROTECTION}}} avoid launch more than 2 x maxParallelThreads to avoid system crash
			for z := 0; z >= 0; z++ { // infinite loop
				time.Sleep(time.Duration(100) * time.Millisecond)
				if(uxmStats.TasksOnWait < maxTasksToLaunchLimit) {
					break;
				}
			}
		}
		uxmStats.TasksOnWait++
		//--
		wg.Add(1)
		go func(LoopId int, i int, TaskId string) {
			defer wg.Done()
			semaphore <- struct{}{} // Lock
			defer func() {
				<-semaphore // Unlock
			}()
			var inListTaskErr bool = false
			uxmStats.TasksProcessed++
			var status = ""
			res, err := http.Get(UrlTaskCall + url.QueryEscape(TaskId)) // execution
			uxmStats.TasksOnWait-- // after execution !
			if err != nil {
				uxmStats.StatALL++
				uxmStats.StatERR++
				fmt.Println(err.Error())
				status = color.RedString("000")
				inListTaskErr = true
			} else {
				res.Body.Close()
				uxmStats.StatALL++
				if(res.StatusCode == 200) {
					uxmStats.Stat200++
					status = color.GreenString(strconv.Itoa(res.StatusCode))
				} else if(res.StatusCode == 202) {
					uxmStats.Stat202++
					status = color.HiCyanString(strconv.Itoa(res.StatusCode))
				} else if(res.StatusCode == 203) {
					uxmStats.Stat203++
					status = color.YellowString(strconv.Itoa(res.StatusCode))
				} else if(res.StatusCode == 208) {
					uxmStats.Stat208++
					status = color.HiMagentaString(strconv.Itoa(res.StatusCode))
				} else if(res.StatusCode == 429) {
					uxmStats.Stat429++
					uxmStats.StatERR++
					inListTaskErr = true
					status = color.HiRedString(strconv.Itoa(res.StatusCode))
				} else if(res.StatusCode == 502) {
					uxmStats.Stat502++
					uxmStats.StatERR++
					inListTaskErr = true
					status = color.HiRedString(strconv.Itoa(res.StatusCode))
				} else if(res.StatusCode == 503) {
					uxmStats.Stat503++
					uxmStats.StatERR++
					inListTaskErr = true
					status = color.HiRedString(strconv.Itoa(res.StatusCode))
				} else if(res.StatusCode == 504) {
					uxmStats.Stat504++
					uxmStats.StatERR++
					inListTaskErr = true
					status = color.HiRedString(strconv.Itoa(res.StatusCode))
				} else {
					uxmStats.StatERR++
					inListTaskErr = true
					status = color.RedString(strconv.Itoa(res.StatusCode))
				}
			}
			if((flagSilent != true) || (inListTaskErr == true)) {
				fmt.Println("Task # " + color.HiYellowString(TaskId) + color.CyanString(" @ Thread.ID:") + color.HiBlackString(strconv.Itoa(LoopId)) + "." + color.HiBlueString(strconv.Itoa(i)) + " :: HTTP Response Status:" + status)
			}
			//-- sleep after each group
			time.Sleep(time.Duration(250) * time.Millisecond)
			//--
		}(LoopId, i, TaskId)
		//-- just add a random pause in milliseconds for give a small breath ... (important for spread of threads in time !!!)
		time.Sleep(time.Duration(rand.Int31n(3) * 25) * time.Millisecond)
		//--
	}
	//--
	if(inListTasks <= 0) {
		uxmStats.EmptyBatches++
	}
	//--
	wg.Wait()
	//--
	fmt.Println("===== Done Processing Batch Tasks [" + strconv.Itoa(inListTasks) + "] =====")
	//--
}


func getIniStrVal(file ini.File, Section string, Key string) string {
	//--
	str, ok := file.Get(Section, Key)
	if !ok {
		str = ""
	}
	//--
	return str
	//--
}


func getIniIntVal(file ini.File, Section string, Key string) int {
	//--
	str := getIniStrVal(file, Section, Key)
	i, err := strconv.Atoi(str)
	if err != nil {
		i = 0
	}
	//--
	return i
	//--
}


func strTrimWhitespaces(s string) string {
	//--
	s = strings.Trim(s, " \t\n\r\x00\x0B")
	//--
	return s
	//--
}


func resetStatistics() {
	//--
	tr := time.Now()
	startTime = tr.Format(time.RFC1123Z)
	//--
	uxmStats.StartTime = startTime
	uxmStats.CurrentTime = startTime
	uxmStats.EmptyBatches = 0
	uxmStats.BatchCycles = 0
	uxmStats.TasksOnWait = 0
	uxmStats.TasksProcessed = 0
	uxmStats.Stat200 = 0
	uxmStats.Stat202 = 0
	uxmStats.Stat203 = 0
	uxmStats.Stat208 = 0
	uxmStats.Stat429 = 0
	uxmStats.Stat502 = 0
	uxmStats.Stat503 = 0
	uxmStats.Stat504 = 0
	uxmStats.StatERR = 0
	uxmStats.StatALL = 0
	//--
}

func uxmHtmlBox(theLogoLink string, theLogoTtlLink string, theHtmlInnerContent string) string {
	//--
	t := time.Now()
	//--
	return `<div style="font-size:2rem; text-align:center;">` + html.EscapeString(t.Format(time.RFC1123Z)) + "</div>" + "\n" + `<div style="background:#778899; color:#FFFFFF; font-size:1.25rem; font-weight:bold; text-align:center; border-radius:3px; padding:10px; margin:20px;">` + "\n" + `<div style="text-align:center; margin:10px; cursor:help;"><a href="` + html.EscapeString(theLogoLink) + `"><img alt="TaskEngine :: ` + html.EscapeString(theLogoTtlLink) + `" title="TaskEngine :: ` + html.EscapeString(theLogoTtlLink) + `" width="64" height="64" src="data:image/svg+xml;base64,` + base64.StdEncoding.EncodeToString([]byte(SVG_LOGO)) + `"></a></div>` + "\n" + theHtmlInnerContent + "\n" + "</div>"
	//--
}


func uxmHtmlSvgSpinner() string {
	//--
	var spinner = SVG_SPIN
	var action = "/stop"
	var txt = "Up and Running :: Click to Pause"
	if(uxmStats.IsActive != true) {
		spinner = strings.Replace(SVG_SPIN, `repeatCount="indefinite"`, `repeatCount="1"`, 1)
		action = "/start"
		txt = "PAUSED :: Click to Activate"
	}
	spinner = `<div style="text-align:center; margin:10px; cursor:help;">` + `<a href="` + html.EscapeString(action) + `"><img fill="freeze" alt="Status: ` + html.EscapeString(txt) + ` ..." title="Status: ` + html.EscapeString(txt) + ` ..." width="64" height="64" src="data:image/svg+xml;base64,` + base64.StdEncoding.EncodeToString([]byte(spinner)) + `"></a>` + `</div>`;
	//--
	return spinner
	//--
}


func main() {

	//--
	testFlagSilent := flag.Bool("s", false, "Silent. Will not output in console all the stuff but only errors.")
	//--
	flag.Parse()
	//--

	//--
	if(*testFlagSilent == true) {
		flagSilent = true
	}
	//--

	//--
	ts := time.Now()
	startTime = ts.Format(time.RFC1123Z)
	//--

	//--
	file, err := ini.LoadFile(INI_FILE)
	if err != nil {
		fmt.Println(color.RedString("ERROR :: INI File cannot be loaded: " + INI_FILE))
		return
	}
	//--
	UrlBatchList = getIniStrVal(file, "URLs", "batch-url")
	if(UrlBatchList == "") {
		fmt.Println(color.RedString("ERROR :: Invalid key [URLs/batch-url] in INI File: " + INI_FILE))
		return
	}
	UrlTaskCall = getIniStrVal(file, "URLs", "tasks-url")
	if(UrlTaskCall == "") {
		fmt.Println(color.RedString("ERROR :: Invalid key [URLs/tasks-url] in INI File: " + INI_FILE))
		return
	}
	parallelWorkers = getIniIntVal(file, "Tunnings", "parallel-workers")
	if((parallelWorkers < 2) || (parallelWorkers > 1024)) {
		parallelWorkers = runtime.NumCPU()
	}
	srvHostAddr = getIniStrVal(file, "MiniServer", "http-addr")
	if(srvHostAddr == "") {
		srvHostAddr = "0.0.0.0"
	}
	srvHostPort = getIniIntVal(file, "MiniServer", "http-port")
	if((srvHostPort < 0) || (srvHostPort > 65535)) {
		srvHostPort = 0
	}
	//--

	//--
	uxmStats = &uxmStuctStats {
		Description: `<TaskEngine GO> Runtime's "Statistics & Settings"`,
		ServiceVersion: THE_VERSION,
		IsActive: true,
		IsSilent: flagSilent,
		StartTime: startTime,
		CurrentTime: startTime,
		BatchURL: UrlBatchList,
		TaskURL: UrlTaskCall,
		NumWorkers: parallelWorkers,
		EmptyBatches: 0,
		BatchCycles: 0,
		TasksOnWait: 0,
		TasksProcessed: 0,
		Stat200: 0,
		Stat202: 0,
		Stat203: 0,
		Stat208: 0,
		Stat429: 0,
		Stat502: 0,
		Stat503: 0,
		Stat504: 0,
		StatERR: 0,
		StatALL: 0,
	}
	//--

	//--
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(color.HiBlackString("\n" + "##### DONE: TaskEngine #####"))
		os.Exit(1)
	}()
	//--
	fmt.Println(color.HiBlackString("\n" + "##### TaskEngine [ " + THE_VERSION + " ] :: Running on #CPUs: " + strconv.Itoa(runtime.NumCPU()) + " #####"))
	//--

	//-- parallel limits
	runtime.GOMAXPROCS(parallelWorkers)
	//-- TLS flexible, allow insecure
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//--

	//--
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//--
		var serverSignature bytes.Buffer
		serverSignature.WriteString("TaskEngine GO MiniServer " + THE_VERSION + "\n")
		serverSignature.WriteString("(c) 2020 unix-world.org" + "\n")
		serverSignature.WriteString("\n")
		serverSignature.WriteString("<Status-URL(s)> :: http://" + srvHostAddr + ":" + strconv.Itoa(srvHostPort) + "/status(.json|.reset)" + "\n")
		serverSignature.WriteString("<Remote-Control-URL(s)> :: " + "http://" + srvHostAddr + ":" + strconv.Itoa(srvHostPort) + "/stop|start|console.verbose|console.silent|quit" + "\n")
		//--
		var statusCode = 202
		//--
		if r.URL.Path != "/" {
			statusCode = 404
			w.WriteHeader(statusCode)
			w.Write([]byte("404 Not Found\n"))
			log.Printf("TaskEngine GO MiniServer :: DEFAULT.NOTFOUND [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
			return
		}
		//--
		log.Printf("TaskEngine GO MiniServer :: DEFAULT [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Refresh", "10")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		fmt.Fprintf(w, HTML_START + uxmHtmlBox("/status", "STATUS", "<pre>" + "\n" + html.EscapeString(serverSignature.String()) + "</pre>") + "\n" + uxmHtmlSvgSpinner() + HTML_END)
		//--
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		//--
		var statusCode = 203
		log.Printf("TaskEngine GO MiniServer :: STATUS [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Refresh", "5")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		t := time.Now()
		uxmStats.CurrentTime = t.Format(time.RFC1123Z)
		var jsonData []byte
		jsonData, err := json.MarshalIndent(uxmStats, "", "\t")
		var jsonStr = "{}"
		if err == nil {
			jsonStr = string(jsonData)
		}
		//--
		fmt.Fprintf(w, HTML_START + uxmHtmlBox("/status.json", "JSON-STATUS", `<div style="font-weight:bold; font-size:1.5rem;">TaskEngine GO MiniServer / Status</div>`) + "\n" + uxmHtmlSvgSpinner() + "\n" + `<div align="center"><div style="text-align:left; width:800px; overflow-x:auto; background:#FAFAFA; border: 1px solid #CCCCCC; border-radius: 5px; padding: 10px;"><pre id="json-code">` + html.EscapeString(jsonStr) + `</pre></div></div>` + "\n" + strTrimWhitespaces(HTML_JS) + HTML_END)
		//--
	})
	http.HandleFunc("/status.json", func(w http.ResponseWriter, r *http.Request) {
		//--
		var statusCode = 200
		log.Printf("TaskEngine GO MiniServer :: STATUS.JSON [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		t := time.Now()
		uxmStats.CurrentTime = t.Format(time.RFC1123Z)
		var jsonData []byte
		jsonData, err := json.Marshal(uxmStats)
		if err == nil {
			w.Write(jsonData)
		} else {
			fmt.Fprintf(w, `{ "ERROR": "Failed to create json structure" }`)
		}
		//--
	})
	http.HandleFunc("/status.reset", func(w http.ResponseWriter, r *http.Request) {
		//--
		flagResetStats = true
		//--
		var statusCode = 208
		log.Printf("TaskEngine GO MiniServer :: STATUS.RESET [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		fmt.Fprintf(w, "OK: Reset Statistics" + "\n")
		//--
	})
	http.HandleFunc("/console.silent", func(w http.ResponseWriter, r *http.Request) {
		//--
		flagSilent = true
		uxmStats.IsSilent = flagSilent
		//--
		var statusCode = 208
		log.Printf("TaskEngine GO MiniServer :: CONSOLE.SILENT [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		fmt.Fprintf(w, "OK: Console Output Silent" + "\n")
		//--
	})
	http.HandleFunc("/console.verbose", func(w http.ResponseWriter, r *http.Request) {
		//--
		flagSilent = false
		uxmStats.IsSilent = flagSilent
		//--
		var statusCode = 208
		log.Printf("TaskEngine GO MiniServer :: CONSOLE.VERBOSE [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		fmt.Fprintf(w, "OK: Console Output Verbose" + "\n")
		//--
	})
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		//--
		uxmStats.IsActive = true
		//--
		var statusCode = 208
		log.Printf("TaskEngine GO MiniServer :: START [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		fmt.Fprintf(w, "OK: Service Active" + "\n")
		//--
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		//--
		uxmStats.IsActive = false
		//--
		var statusCode = 208
		log.Printf("TaskEngine GO MiniServer :: STOP [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		fmt.Fprintf(w, "OK: Service Paused" + "\n")
		//--
	})
	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		//--
		var statusCode = 208
		log.Printf("TaskEngine GO MiniServer :: QUIT [%s %s %s] %s [%s] %s\n", r.Method, r.URL, r.Proto, strconv.Itoa(statusCode), r.Host, r.RemoteAddr)
		//--
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(statusCode) // status code must be after content type
		//--
		fmt.Fprintf(w, "OK: Service will STOP and EXIT in 10 seconds" + "\n")
		//--
		go func() {
			time.Sleep(time.Duration(10) * time.Second)
			fmt.Println(color.HiYellowString("INFO: Service will QUIT as requested by Remote Control via HTTP"))
			os.Exit(0)
		}()
		//--
	})
	//--
	ln, err := net.Listen("tcp", srvHostAddr + ":" + strconv.Itoa(srvHostPort))
	if err != nil {
		log.Fatal(err)
	}
	bindTcpAddr = ln.Addr().String()
	fmt.Println("HTTP Built-in MiniServer Listening on TCP: " + bindTcpAddr)
	//--
	go func() {
		log.Fatal(http.Serve(ln, nil))
	}()
	//--

	//-- main loop for task engine
	for i := 0; i >= 0; i++ { // infinite loop
		//--
		tasks(parallelWorkers, i)
		time.Sleep(time.Duration(rand.Int31n(2)) * time.Second)
		//--
	}
	//--

}

// #END
