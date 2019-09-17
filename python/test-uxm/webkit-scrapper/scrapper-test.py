#!/usr/bin/env python2.7
# -*- coding: utf-8 -*-
#
# copyright (c) 2017 Francesco Guarnieri
#
# copyright (c) 2017 unix-world.org
#

import os
from os import path
import sys
import time

import gi
gi.require_version('Gtk', '3.0')
gi.require_version('WebKit2', '4.0')
from gi.repository import Gtk, Gdk, Gio, Pango, WebKit2

start_url = "https://play.google.com/store/apps/category/COMMUNICATION/collection/topselling_free?hl=en"

# A script stack to execute after each page is loaded
scripts = [
"""
	window.scrollTo(0,document.body.scrollHeight);
	function myFunction() {
		var mylinks = [];
		jQuery('a.title').each(function(index) {
			mylinks.push(jQuery(this).attr('href'));
		});
		var ddata = JSON.stringify(mylinks);
		mylinks = [];
		var area_id = 'wk-py-uxm-data';
		jQuery('body').append('<textarea id="' + area_id + '"></textarea>');
		jQuery('textarea#' + area_id).val(ddata).select();
		ddata = '';
		document.execCommand("Copy");
		jQuery('textarea#' + area_id).remove();
		try {
			window.webkit.messageHandlers.COMM_CHANNEL_JS_WK.postMessage('UxmSignalOne');
		} catch(err){}
	}
	setTimeout(function(){ myFunction(); }, 3000);
"""
]

class Browser(Gtk.Window):
    def __init__(self):

        # security: clear the clipboard because the Javascript can access it and avoid access sensitive data
        clipboard = Gtk.Clipboard.get(Gdk.SELECTION_CLIPBOARD)
        clipboard.clear()

        Gtk.Window.__init__(self)
        self.set_position(Gtk.WindowPosition.CENTER)
        self.set_default_size(1152, 800)

        self.scripts = None
        contentManager = WebKit2.UserContentManager()
        contentManager.connect("script-message-received::COMM_CHANNEL_JS_WK", self.__handleScriptMessage)
        if not contentManager.register_script_message_handler("COMM_CHANNEL_JS_WK"):
            print("Error registering script message handler: COMM_CHANNEL_JS_WK")

        # Init the WebView with the contentManager
        mydir = path.abspath(path.dirname(__file__))
        print("Extension directory:", mydir)
        ctx = WebKit2.WebContext.get_default()
        ctx.set_web_extensions_directory(mydir)
        ctx.set_tls_errors_policy(WebKit2.TLSErrorsPolicy.IGNORE) # ignore TLS Errors
        self.web_view = WebKit2.WebView.new_with_context(ctx).new_with_user_content_manager(contentManager)
        self.web_view.connect("load-changed", self.__loadFinishedCallback)

        # Settings: Javascript messages need to be accessed via clipboard
        settings = WebKit2.Settings()
        settings.set_property('javascript-can-access-clipboard', True)
        # Apply the rest of settings
        settings.set_property('user-agent', 'iPad ewhjwhdsej')
        settings.set_property('enable-plugins', False)
        settings.set_property('enable-java', False)
        settings.set_property('enable-javascript', True)
        settings.set_property('allow-modal-dialogs', True)
        settings.set_property('enable-private-browsing', True)
        settings.set_property('enable-page-cache', True)
        settings.set_property('enable-smooth-scrolling', True)
        settings.set_property('enable-write-console-messages-to-stdout', True)
        self.web_view.set_settings(settings)

        cancelButton = Gtk.Button()
        icon = Gio.ThemedIcon(name="window-close-symbolic")
        image = Gtk.Image.new_from_gicon(icon, Gtk.IconSize.SMALL_TOOLBAR)
        cancelButton.add(image)
        cancelButton.connect("clicked", self.__close)

        headerBar = Gtk.HeaderBar()
        headerBar.set_show_close_button(False)
        self.set_titlebar(headerBar)
        boxTitle = Gtk.Box(spacing = 6)
        self.spinner = Gtk.Spinner()
        boxTitle.add(self.spinner)
        labelTitle = Gtk.Label()
        labelFont = labelTitle.get_pango_context().get_font_description()
        labelFont.set_weight(Pango.Weight.BOLD)
        labelTitle.modify_font(labelFont)
        labelTitle.set_text(start_url)
        boxTitle.add(labelTitle)
        headerBar.set_custom_title(boxTitle)

        self.stopButton = Gtk.Button()
        icon = Gio.ThemedIcon(name="media-playback-stop-symbolic")
        image = Gtk.Image.new_from_gicon(icon, Gtk.IconSize.SMALL_TOOLBAR)
        self.stopButton.add(image)
        self.stopButton.connect("clicked", self.__on_stop_click)
        headerBar.pack_end(self.stopButton)

        box = Gtk.Box()
        Gtk.StyleContext.add_class(box.get_style_context(), "linked")
        box.add(cancelButton)
        headerBar.pack_start(box)
        browserBox = Gtk.Box()
        browserBox.set_orientation(Gtk.Orientation.VERTICAL)
        browserBox.pack_start(self.web_view, True, True, 0)
        self.add(browserBox)

    # Callback called up at each change of state of the load phase of the webview
    # Possible Events: WEBKIT_LOAD_STARTED, WEBKIT_LOAD_REDIRECTED, WEBKIT_LOAD_COMMITTED, WEBKIT_LOAD_FINISHED
    def __loadFinishedCallback(self, web_view, load_event):
        if self.scripts and (len(self.scripts) > 0) and (load_event == WebKit2.LoadEvent.FINISHED):
            self.__waitMode(False)
            web_view.run_javascript_from_gresource("https://code.jquery.com/jquery-3.2.1.min.js")
            time.sleep(2)
            web_view.run_javascript(self.scripts.pop(), None, self.__javascript_finished, None)
            return False

    # Callback to manage the end of execution of Javascript Code
    def __javascript_finished(self, webview, task, user_data = None):
        try:
            # the js_result object of JavaScriptResult type is not manageable in Python
            result = webview.run_javascript_finish(task)
        except Exception as e:
            print("JAVASCRIPT ERROR MSG: %s" % e)

    # Handle the message from Javascript, async mode (not necessary to end the execution)
    def __handleScriptMessage(self, contentManager, js_result):
        # Here the js_result object of JavaScriptResult type is not manageable in Python, thus will use GDK Clipboard to transfer the data
        print "[ Data received from JavaScript ]"
        #context = WebKit2.JavascriptResult.get_global_context(js_result)
        #print context
        #uri = self.web_view.get_uri()
        #print uri
        clipboard = Gtk.Clipboard.get(Gdk.SELECTION_CLIPBOARD)
        resultStr = clipboard.wait_for_text()
        if resultStr:
            fh = open("py-scrap.json","w")
            fh.write("%s" % resultStr)
            fh.close()
            #print "%s" % resultStr
            clipboard.clear()
            print("[ Cleaning Up ... Done ]")
            time.sleep(3)
            sys.exit()
#           os.execl(sys.executable, sys.executable, *sys.argv)
            return True

    # Set the wait mode if it is active or not
    def __waitMode(self, toggle):
        self.stopButton.set_sensitive(toggle)
        self.stopButton.set_visible(toggle)
        self.web_view.set_sensitive(not toggle)
        if toggle: self.spinner.start()
        else: self.spinner.stop()

    # Manages the click of the stop button
    def __on_stop_click(self, widget):
        self.__waitMode(False)
        self.web_view.stop_loading()

    # Manage the close of the mini-browser
    def __close(self, widget=None):
        self.web_view.stop_loading()
        self.destroy()
        Gtk.main_quit()

    # Opens the site
    def open(self, site, scripts = None):
        self.scripts = scripts
        self.__waitMode(True)
        self.web_view.load_uri(site)

if __name__ == "__main__":
    browser = Browser()
    browser.show_all()
    if len(sys.argv) > 1:
        start_url = sys.argv[1]
    browser.open(start_url, scripts)
    Gtk.main()

#END
