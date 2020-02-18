#!/usr/bin/env python2.7
# -*- coding: utf-8 -*-

# Webkit Crawler
# (c) 2017-2020 Radu.I

uxmScriptVersion = 'v.20200213.1203'

import yaml, random
with open("crawler.yaml", 'r') as ymlfile:
    cfgdata = yaml.safe_load(ymlfile)
#   print(cfgdata)

userAgent = cfgdata['userAgent']
urlScript = cfgdata['urlScript']
pathjQuery = cfgdata['pathjQuery']
pathjScript = cfgdata['pathjScript']

import sys
import os
from os import path
import time

import gi
gi.require_version('Gtk', '3.0')
#gi.require_version('Gdk', '3.0')
gi.require_version('WebKit2', '4.0')
from gi.repository import WebKit2, Gtk #, Gdk

#mydir = path.abspath(path.dirname(__file__))
mydir = path.abspath('./webkit-py-extensions')
#print("Current Webkit Extension directory:", mydir)

#GObject.threads_init()

ctx = WebKit2.WebContext.get_default()
ctx.set_web_extensions_directory(mydir)
ctx.set_tls_errors_policy(WebKit2.TLSErrorsPolicy.IGNORE) # ignore TLS Errors
#ctx.set_web_extensions_initialization_user_data(GLib.Variant.new_string("test string"))

contentManager = WebKit2.UserContentManager()

web = WebKit2.WebView.new_with_context(ctx).new_with_user_content_manager(contentManager)
cfgweb = WebKit2.Settings()
cfgweb.set_property('default-charset', 'UTF-8')
cfgweb.set_property('user-agent', userAgent)
cfgweb.set_property('enable-plugins', False)
cfgweb.set_property('enable-java', False)
cfgweb.set_property('enable-javascript', True)
cfgweb.set_property('javascript-can-access-clipboard', False)
cfgweb.set_property('javascript-can-open-windows-automatically', False)
cfgweb.set_property('allow-modal-dialogs', True)
cfgweb.set_property('enable-private-browsing', True)
cfgweb.set_property('enable-page-cache', True)
cfgweb.set_property('enable-smooth-scrolling', False)
cfgweb.set_property('enable-write-console-messages-to-stdout', True)

cfgweb.set_property('enable-webaudio', False)
cfgweb.set_property('enable-webgl', False)
cfgweb.set_property('enable-accelerated-2d-canvas', False)
cfgweb.set_property('hardware-acceleration-policy', WebKit2.HardwareAccelerationPolicy.NEVER)

web.set_settings(cfgweb)

#wkinfo = WebKit2.ApplicationInfo;
uxmScriptSignature = 'Python v.' + str(sys.version_info[0]) + '.' + str(sys.version_info[1]) + '.' + str(sys.version_info[2]) + ' / WebkitGtk v.' + str(WebKit2.get_major_version()) + '.' + str(WebKit2.get_minor_version()) + '.' + str(WebKit2.get_micro_version()) + ' :: Crawler ' + uxmScriptVersion
print uxmScriptSignature

script = """
var uxmScriptUrl = '""" + urlScript + """';
""";
with open (pathjQuery, "r") as myfile1:
    script += "\n" + myfile1.read()
with open (pathjScript, "r") as myfile2:
    script += "\n" + myfile2.read()
#print script

wnd = Gtk.Window()
wnd.set_resizable(False)
wnd.connect("destroy", Gtk.main_quit)

## Fix for crashing on X: https://stackoverflow.com/questions/14277724/python-gtk-webkit-with-css-issue-in-width
box = Gtk.VBox(homogeneous=False, spacing=0)
wnd.add(box)
box.pack_start(web, expand=False, fill=True, padding=0)
web.set_size_request(1280, 700)
##
#wnd.add(web)
#wnd.set_default_size(1152, 800)
## end Fix

addressbar = Gtk.Entry()
addressbar.set_editable(False)
addressbar.set_can_focus(False)
box.pack_start(addressbar, expand=False, fill=False, padding=0)

progressbar = Gtk.ProgressBar()
progressbar.set_show_text(False)
box.pack_start(progressbar, expand=False, fill=False, padding=0)

wnd.show_all()

def uxm_rand_sleep_timer():
    return random.randint(3, 7)

randTest = uxm_rand_sleep_timer()
#print('RandTest = ' + str(randTest))

def js_exec_finished(webview, task, user_data = None):
    try:
        # the js_result object of JavaScriptResult type is not manageable in Python
        result = webview.run_javascript_finish(task)
    except Exception as e:
        errMsg = "%s" % e
        if errMsg:
            print("JAVASCRIPT ERROR MSG: ", errMsg)

def on_load_changed(webview, event):
    crrURL = webview.get_uri()
    if event == WebKit2.LoadEvent.FINISHED: # Possible Events: WEBKIT_LOAD_STARTED, WEBKIT_LOAD_REDIRECTED, WEBKIT_LOAD_COMMITTED, WEBKIT_LOAD_FINISHED
        if crrURL.startswith('data:'):
            arrDatUrl = crrURL.split(",")
            DatUrlType = str(arrDatUrl[0])
            DatUrlUri = str(arrDatUrl[1].decode('base64'))
            if(DatUrlType == 'data:text/plain;base64' and DatUrlUri != ''):
                wndTtl = "DataURL: " + DatUrlUri
                wnd.set_title(wndTtl[0:100])
                time.sleep(uxm_rand_sleep_timer())
                web.load_uri(DatUrlUri)
            else:
                print("Invalid Data: ", crrURL)
        else:
            wndTtl = webview.get_title()
            try:
                wnd.set_title(wndTtl[0:100])
            except:
                print("Could not set page Title ...")
            addressbar.set_text(crrURL)
            progressbar.set_fraction(1)
            time.sleep(uxm_rand_sleep_timer())
            webview.run_javascript(script, None, js_exec_finished, None)
    else:
        wnd.set_title("Loading URL ...  {:0.1f}%".format(webview.get_estimated_load_progress() * 100))
        if crrURL.startswith('data:'):
            addressbar.set_text("Processing Data ...")
        else:
            addressbar.set_text(crrURL)
        progressbar.set_fraction(webview.get_estimated_load_progress())

def on_load_tlserrors(webview, url, certificate, error):
    print("Error TLS: ", url, error, certificate)
    return True

def on_load_failed(webview, event, url, error):
    print("Error loading", url, "-", error)

def on_js_message(contentManager, js_result):
    print("[ Webkit Js Signal Handler: QUIT ]")
    time.sleep(1)
    sys.exit()

contentManager.connect("script-message-received::COMM_CHANNEL_JS_WK", on_js_message)
if not contentManager.register_script_message_handler("COMM_CHANNEL_JS_WK"):
    print("Error registering script message handler: COMM_CHANNEL_JS_WK")

web.connect("load-changed", on_load_changed)
web.connect("load-failed-with-tls-errors", on_load_tlserrors)
web.connect("load-failed", on_load_failed)

#if len(sys.argv) > 1:
#    web.load_uri(sys.argv[1])
#else:
web.load_uri(urlScript)

Gtk.main()

#END
