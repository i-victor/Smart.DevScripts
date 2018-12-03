
// GO Lang
// go build test-webview.go (on openbsd may need to: CGO_LDFLAGS_ALLOW='-Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib' go build test-webview.go)
// (c) 2017-2018 unix-world.org
// version: 2018.12.02

package main

import (
	"github.com/zserge/webview"
//	"github.com/unix-world/webview"
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

// #END
