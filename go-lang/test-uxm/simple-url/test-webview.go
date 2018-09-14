
package main

import (
	"github.com/zserge/webview"
)

func main() {
	w := webview.New(webview.Settings{
		Width:  960,
		Height: 720,
		Title:  "Test Loading External URL",
		URL:    "https://demo.unix-world.org/smart-framework/",
	})
	defer w.Exit()
	w.Run()
}
