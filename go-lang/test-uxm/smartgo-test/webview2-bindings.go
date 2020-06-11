
// GO Lang
// go build test-webview.go (on openbsd may need to: CGO_LDFLAGS_ALLOW='-Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib' go build webview2-bindings.go)
// (c) 2020 unix-world.org
// version: 20200517

package main

/*
#cgo openbsd LDFLAGS: -Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib
*/

import (
	"log"
//	webview "github.com/zserge/webview2"
	webview "github.com/unix-world/smartgo/webview2"
)

func main() {

	debug := false
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)


	w.Bind("noop", func() string {
		log.Println("run: noop")
		return "noop"
	})

	w.Bind("add", func(a int, b int) int {
		var c int = a + b
		log.Println("run: add (", c , ")")
		return c
	})

	w.Bind("quit", func() {
		w.Terminate()
	})


//	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")


	w.Navigate(`data:text/html,<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Hello World</title>
				<script>
				function clickButton1() {
					noop().then(function(res) {
						console.log('noop res', res);
						add(1, 2).then(function(res) {
							console.log('add res', res);
						});
					});
				}
				function clickButton2() {
					quit().then(function(){});
				}
				</script>
			</head>
			<body>
				<div id="hello"></div>
				<button onClick="clickButton1()">Run Action Noop + Add</button> <button onClick="clickButton2()">Quit</button>
				<script>
					document.getElementById('hello').innerText = ` + "`hello, ${navigator.userAgent}`" + `;
				</script>
			</body>
		</html>`)


	w.Run()

}

// #END
