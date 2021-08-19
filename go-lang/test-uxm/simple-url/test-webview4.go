
// GO Lang
// go build test-webview4.go (on openbsd may need to: CGO_LDFLAGS_ALLOW='-Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib' go build test-webview4.go)
// (c) 2017-2021 unix-world.org
// version: 2021.06.19

package main

/*
#cgo openbsd LDFLAGS: -Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib
*/

import (
	"github.com/webview/webview"
)

func main() {
    debug := false
    w := webview.New(debug)
    defer w.Destroy()
    w.SetTitle("Minimal webview example")
    w.SetSize(960, 720, webview.HintNone)
    w.Navigate("http://demo.unix-world.org/smart-framework/")
    w.Run()
}

// #END
