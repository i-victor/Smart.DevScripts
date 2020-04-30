
// GO Lang
// go build test-webview.go (on openbsd may need to: CGO_LDFLAGS_ALLOW='-Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib' go build test-webview.go)
// (c) 2017-2018 unix-world.org
// version: 2018.12.02

package main

/*
#cgo openbsd LDFLAGS: -Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib
*/

import (
	"log"
	webview "github.com/zserge/webview2"
)

func main() {

	debug := false
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)


	w.Bind("noop", func() string {
		log.Println("hello")
		return "hello"
	})

/*
	w.Bind("add", func(a, b int) int {
		return a + b
	})
	w.Bind("quit", func() {
		w.Terminate()
	})
*/

//	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")


	w.Navigate(`data:text/html,
		<!doctype html>
		<html>
			<body>hello</body>
			<script>
				window.onload = function() {
					document.body.innerText = ` + "`hello, ${navigator.userAgent}`" + `;
					noop().then(function(res) {
						console.log('noop res', res);
						add(1, 2).then(function(res) {
							console.log('add res', res);
							quit();
						});
					});
				};
			</script>
		</html>
	)`)


	w.Run()

}

// #END
