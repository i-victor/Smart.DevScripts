
// GoLang: GTK+3 / WebKit

package main

// #cgo openbsd LDFLAGS:-Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib
import (
	"log"
	"fmt"
	"runtime"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	wk2gtk3 "github.com/unix-world/smartgo/webkit2gtk3"
)


func main() {

	runtime.LockOSThread()

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")

	wc := wk2gtk3.DefaultWebContext()
//	wc := wk2gtk3.PrivateBrowsingWebContext()
	wc.ClearCache()
	wc.TlsPolicyIgnoreErrors() // security (does not work with PrivateBrowsingWebContext)

	webView := wk2gtk3.NewWebView()
	ws := webView.Settings()
	ws.EnableFullScreen(false)
	ws.SetDefaultCharset("UTF-8")
	ws.JavascriptCanAccessClipboard(false) // security !
	ws.JavascriptCanOpenWindowsAutomatically(false) // security
	ws.EnableXssAuditor(true) // security
	ws.EnableDnsPrefetching(false)
//	ws.AllowDataUrls(true) // allow data URLs, this is default ... ; v 2.28+ only
	ws.EnableJava(false)
	ws.EnableJavascript(true)
	ws.EnablePlugins(false)
	ws.EnablePageCache(false)
	ws.EnableSmoothScrolling(false)
	ws.EnableMedia(false)
	ws.EnableWebAudio(false)
	ws.EnableWebGl(false)
	ws.EnableAccelerated2DCanvas(false)
	ws.SetHardwareAccelerationPolicy("WEBKIT_HARDWARE_ACCELERATION_POLICY_NEVER")
//	ws.SetCustomUserAgent("Test Browser")
	ws.SetUserAgentWithApplicationDetails("GoLang WebKit Test", "1.0")

//	ws.SetEnableWriteConsoleMessagesToStdout(true) // debug only
//	ws.EnableDeveloperExtras(true)

	win.Connect("destroy", func() {
		webView.Destroy()
		gtk.MainQuit()
	})

	var loadFailed bool = false
	webView.Connect("load-failed", func() {
		loadFailed = true
		fmt.Println("Load Failed !") // malformed URL, network is down, ...
	})
	webView.Connect("load-failed-with-tls-errors", func() {
		loadFailed = true
		fmt.Println("Load Failed with TLS Errors ...") // on https the TLS certificate is invalid, expired, handshake failed, ...
	})

	var theJavaScript string = ""
//	theJavaScript = "jQuery('div').html();"
	theJavaScript = "(function() { alert('Hello World, a message from GoLang to WebKit !'); return true; })();"
	webView.Connect("load-changed", func(_ *glib.Object, i int) {
		//--
		if(loadFailed == true) {
			return // do not execute this function if load failed !!
		}
		//--
		loadEvent := wk2gtk3.LoadEvent(i)
		//--
		switch loadEvent {
			case wk2gtk3.LoadFinished:
				//--
				fmt.Println("Load Finished")
				//--
				if(loadFailed == false) {
					//--
					fmt.Printf("Title: %q\n", webView.Title())
					fmt.Printf("URI: %s\n", webView.URI())
					//--
					if(theJavaScript != "") {
						webView.RunJavaScript(theJavaScript, func(status string, result string) {
							if(status != "OK") {
								fmt.Println("JavaScript Failed: " + status + ": " + result)
							} else {
								fmt.Println("JavaScript: " + status + ": " + result)
							}
						})
					}
					//--
				}
				//--
				break;
			case wk2gtk3.LoadStarted:
				fmt.Println("Load Started")
				break;
			case wk2gtk3.LoadCommitted:
				fmt.Println("Load Committed")
				break;
			case wk2gtk3.LoadRedirected:
				fmt.Println("Load Redirected")
				break;
			default:
				fmt.Println("Load Event: ", loadEvent)
		} //end switch
		//---
	})

//	glib.IdleAdd(func() bool {
		webView.LoadURI("https://demo.unix-world.org/smart-framework/")
//		return false
//	})

// Add the label to the window.
	win.Add(webView)

	// Set the default window size.
	win.SetDefaultSize(960, 720)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	gtk.Main()

}


// #END
