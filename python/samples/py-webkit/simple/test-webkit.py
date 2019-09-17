#!/usr/bin/env python2.7
# -*- coding: utf-8 -*-

import sys
from os import path
import time

import gi
gi.require_version('Gtk', '3.0')
#gi.require_version('Gdk', '3.0')
gi.require_version('WebKit2', '4.0')
from gi.repository import WebKit2, Gtk #, Gdk

mydir = path.abspath(path.dirname(__file__))
print("Extension directory:", mydir)

ctx = WebKit2.WebContext.get_default()
ctx.set_web_extensions_directory(mydir)
ctx.set_tls_errors_policy(WebKit2.TLSErrorsPolicy.IGNORE) # ignore TLS Errors
#ctx.set_web_extensions_initialization_user_data(GLib.Variant.new_string("test string"))

web = WebKit2.WebView.new_with_context(ctx)
cfgweb = WebKit2.Settings()
#cfgweb.set_property('user-agent', 'iPad')
cfgweb.set_property('enable-javascript', True)
web.set_settings(cfgweb)


wnd = Gtk.Window()
wnd.connect("destroy", Gtk.main_quit)
wnd.add(web)
wnd.set_default_size(1152, 800)
wnd.show_all()

def on_load_changed(webview, event):
    if event == WebKit2.LoadEvent.FINISHED:
        try:
            wnd.set_title(webview.get_title())
        except:
            print("Could not set page Title ...")
#       time.sleep(5)
#       wnd.resize(1280, 960)
#       web.run_javascript_from_gresource("https://code.jquery.com/jquery-3.2.1.min.js")
#       web.run_javascript("$('html, body').animate({ scrollTop: $(document).height() }, 1000);")
        time.sleep(1)
        web.run_javascript("$('html, body').animate({ scrollTop: $(document).height() }, 2000);")
#       time.sleep(5)
#       web.run_javascript("var mylinks = []; $('a.title').each(function(index) { mylinks.push($(this).attr('href')); }); $('body').empty().html(mylinks.join('<br>'))")
    else:
        wnd.set_title("Loading ...  {:0.1f}%".format(webview.get_estimated_load_progress()))

def on_load_tlserrors(webview, url, certificate, error):
    print("Error TLS: ", url, error, certificate)
    return True

def on_load_failed(webview, event, url, error):
    print("Error loading", url, "-", error)

web.connect("load-changed", on_load_changed)
web.connect("load-failed-with-tls-errors", on_load_tlserrors)
web.connect("load-failed", on_load_failed)

if len(sys.argv) > 1:
    web.load_uri(sys.argv[1])
else:
    web.load_uri("https://127.0.0.1/sites/tests/test-ua.php")

Gtk.main()
