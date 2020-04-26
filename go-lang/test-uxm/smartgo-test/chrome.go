
// GO Lang

package main

import (
	"os/exec"
	"runtime"
	"strconv"
)

// openBrowser tries to open the URL in a browser,
// and returns whether it succeed in doing so.
// Windows 10 fails to open URLs containing ampersands(&) because they are interpreted as a command separator by cmd.exe,
// they need to be escaped by a caret(^) i.e. strings.Replace(url, "&", "^&", -1)
func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}


func openIridiumAppBrowser(url string, sizeW int, sizeH int) bool {
	if(sizeW < 600) {
		sizeW = 600
	} else if(sizeW > 1600) {
		sizeW = 1600
	}
	if(sizeH < 400) {
		sizeH = 400
	} else if(sizeH > 1080) {
		sizeH = 1080
	}
	cmd := exec.Command("iridium", "--app=" + url, "--window-position=100,100", "--window-size=" + strconv.Itoa(sizeW) + "," + strconv.Itoa(sizeH))
	return cmd.Start() == nil
}


func main() {

//	openBrowser("https://127.0.0.1/sites/")
	openIridiumAppBrowser("https://127.0.0.1/sites/", 600, 400)

}

// #END
