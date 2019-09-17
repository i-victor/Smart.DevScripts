#!/usr/bin/env python2.7
# -*- coding:utf-8 -*-
#
# Copyright (C) 2013 Carlos Jenkins <carlos@jenkins.co.cr>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

import gi
gi.require_version('Gtk', '3.0')
gi.require_version('WebKit2', '4.0')
from gi.repository import Gtk, WebKit2  # gir1.2-webkit2-3.0
from os.path import abspath, dirname, join

WHERE_AM_I = abspath(dirname(__file__))


local_uri = 'webbrowser://'
initial_html = '''\
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>WebView Recipe</title>
</head>

<body>
<p>Some links:</a>
<ul>
    <li><a href="https://www.google.com/">Google</a></li>
    <li><a href="http://www.gtk.org/">Gtk+</a></li>
    <li><a href="https://glade.gnome.org/">Glade</a></li>
    <li><a href="http://www.python.org/">Python</a></li>
    <li><a href="file:///etc/hosts">Local <tt>/etc/hosts</tt> file</a></li>
</ul>
</body>
</html>
'''


class WebBrowser(object):
    """
    Simple WebBrowser class.

    Uses WebKit introspected bindings for Python:
    http://webkitgtk.org/reference/webkit2gtk/stable/
    """

    def __init__(self):
        """
        Build GUI
        """
        # Build GUI from Glade file
        self.builder = Gtk.Builder()
        self.builder.add_from_file(join(WHERE_AM_I, 'webbrowser.ui'))

        # Get objects
        go = self.builder.get_object
        self.window = go('window')
        self.entry = go('entry')
        self.scrolled = go('scrolled')

        # WebView Context
        ctx = WebKit2.WebContext.get_default()
#       ctx.set_tls_errors_policy(WebKit2.TLSErrorsPolicy.IGNORE) # ignore TLS Errors

        # Create WebView
        self.webview = WebKit2.WebView().new_with_context(ctx)
        self.scrolled.add_with_viewport(self.webview)

        # Connect signals
        self.builder.connect_signals(self)
        self.webview.connect('load-changed', self.load_changed_cb)
        self.window.connect('delete-event', Gtk.main_quit)

        # Everything is ready
        self.load_uri(local_uri + 'home')
        self.window.show_all()

    def entry_cb(self, widget, user_data=None):
        """
        Callback when Enter is pressed in URL entry.
        """
        self.load_uri(self.entry.get_text())

    def load_changed_cb(self, webview, event, user_data=None):
        """
        Callback for when the load operation in webview changes.
        """
        ev = str(event)
        if 'WEBKIT_LOAD_COMMITTED' in ev:
            self.entry.set_text(self.webview.get_uri())

    def load_uri(self, uri):
        """
        Load an URI on the browser.
        """
        self.entry.set_text(uri)
        if uri.startswith(local_uri):
            section = uri[len(local_uri):]
            if section == 'home':
                self.webview.load_html(initial_html, local_uri)
                return
        self.webview.load_uri(uri)
        return


if __name__ == '__main__':
    gui = WebBrowser()
    Gtk.main()
