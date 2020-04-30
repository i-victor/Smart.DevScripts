
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

	//--
	runtime.LockOSThread()
	//--
	gtk.Init(nil)
	//--

	//--
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create GTK Window:", err)
	} //end if
	win.SetTitle("Simple Webkit2 GTK3 Example")
	//--

	//-- complex widgets
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Unable to create GTK VBox:", err)
	} //end if
	vbox.Set("homogeneous", false)
	//--
	titlebar, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create GTK Entry TitleBar:", err)
	} //end if
	titlebar.SetEditable(false)
	titlebar.SetCanFocus(false)
	//--
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal("Unable to create GTK HBox:", err)
	} //end if
	box.Set("homogeneous", false)
	//--
	btn, err := gtk.ButtonNewWithLabel("Home")
//	btn, err := gtk.ButtonNewWithIcon("gtk-home") // this is a custom button type implemented by unixman
	if err != nil {
		log.Fatal("Unable to create GTK Button:", err)
	} //end if
	//--
	addressbar, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create GTK Entry AddressBar:", err)
	} //end if
	addressbar.SetEditable(false)
	addressbar.SetCanFocus(false)
	//--
	progressbar, err := gtk.ProgressBarNew()
	if err != nil {
		log.Fatal("Unable to create GTK ProgressBar:", err)
	} //end if
	progressbar.SetShowText(false)
	//--

	//--
	var homeURL string = "https://demo.unix-world.org/smart-framework/"
	var crrURL string = homeURL
	//--
	wc := wk2gtk3.DefaultWebContext()
//	wc := wk2gtk3.PrivateBrowsingWebContext()
	wc.ClearCache()
	wc.TlsPolicyIgnoreErrors() // security (does not work with PrivateBrowsingWebContext)
	//--
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
	//--
//	ws.SetEnableWriteConsoleMessagesToStdout(true) // debug only
//	ws.EnableDeveloperExtras(true)
	//--
	win.Connect("destroy", func() {
		webView.Destroy()
		gtk.MainQuit()
	})
	//--
	var loadFailed bool = false
	webView.Connect("load-failed", func() {
		loadFailed = true
		fmt.Println("Load Failed !") // malformed URL, network is down, ...
	})
	webView.Connect("load-failed-with-tls-errors", func() {
		loadFailed = true
		fmt.Println("Load Failed with TLS Errors ...") // on https the TLS certificate is invalid, expired, handshake failed, ...
	})
	//--

	//--
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
			//	fmt.Println("Load Finished")
				progressbar.SetFraction(1)
				//--
				if(loadFailed == false) {
					//--
					crrURL = webView.URI()
					addressbar.SetText(crrURL)
					titlebar.SetText(webView.Title())
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
			//	fmt.Println("Load Started")
				crrURL = webView.URI()
				addressbar.SetText(crrURL)
				titlebar.SetText("< ... >")
				progressbar.SetFraction(0.1)
				break;
			case wk2gtk3.LoadCommitted:
			//	fmt.Println("Load Committed")
				titlebar.SetText("< ... Loading ... >")
				progressbar.SetFraction(0.5)
				break;
			case wk2gtk3.LoadRedirected:
			//	fmt.Println("Load Redirected")
				titlebar.SetText("< ... Redirecting ... >")
				progressbar.SetFraction(0.25)
				break;
			default:
				fmt.Println("Load Event: ", loadEvent)
		} //end switch
		//---
	})
	//--

	//--
//	glib.IdleAdd(func() bool { // used to use webview in background
	webView.LoadURI(crrURL)
//	return false
//	})

	//-- simple
	// Add the label to the window.
//	win.Add(webView)
	// Set the default window size.
//	win.SetDefaultSize(960, 720)
	//-- (alternate to simple) complex, within a vbox
	webView.SetSizeRequest(960, 720)
	//--
	win.Add(vbox)
	//--
	vbox.PackStart(titlebar, false, false, 0) // obj, expand, fill, padding
	//--
	vbox.PackStart(webView, false, true, 1) // obj, expand, fill, padding
	//--
	btn.Connect("released", func(_ *gtk.Button) {
		wc.ClearCache()
		webView.LoadURI(homeURL)
	})

	box.PackStart(btn, false, true, 0) // obj, expand, fill, padding
	//--
	addressbar.SetSizeRequest(900, 10)
	box.PackStart(addressbar, true, true, 0) // obj, expand, fill, padding
	//--
	vbox.PackStart(box, false, true, 0) // obj, expand, fill, padding
	//--
	vbox.PackStart(progressbar, false, false, 0) // obj, expand, fill, padding
	//--


	//-- Recursively show all widgets contained in this window.
	win.ShowAll()
	//--
	gtk.Main() // main loop
	//--

}


// #END
