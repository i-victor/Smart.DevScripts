#!/usr/bin/env python2.7
# -*- coding: utf-8 -*-

# Webkit HTTPS Server with CORS
# (c) 2017-2019 Radu.I

uxmScriptVersion = 'v.20190710.1159'

import BaseHTTPServer, SimpleHTTPServer
import ssl
import os, sys, signal

if sys.argv[1:]:
    port = int(sys.argv[1])
else:
    port = 4443

web_dir = os.path.join(os.path.dirname(__file__), 'web')
os.chdir(web_dir)

class CORSHTTPRequestHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_my_headers()

        SimpleHTTPServer.SimpleHTTPRequestHandler.end_headers(self)

    def send_my_headers(self):
        self.send_header("Access-Control-Allow-Origin", "*")

requestHandler = CORSHTTPRequestHandler
#requestHandler = SimpleHTTPServer.SimpleHTTPRequestHandler
httpd = BaseHTTPServer.HTTPServer(('localhost', port), requestHandler)
httpd.socket = ssl.wrap_socket(httpd.socket, certfile='../cert.pem', server_side=True)
#httpd.serve_forever()

# A custom signal handle to allow us to Ctrl-C out of the process
def signal_handler(signal, frame):
    print('Exiting HTTPS server (Ctrl+C pressed)')
    try:
      if(httpd):
        httpd.server_close()
    finally:
      sys.exit(0)

# Install the keyboard interrupt handler
signal.signal(signal.SIGINT, signal_handler)

# Now loop forever
try:
    print('Starting HTTPS Server.Py (' + uxmScriptVersion + ') on localhost:' + str(port) + ' (press Ctrl+C to Stop ...)')
    while True:
        sys.stdout.flush()
        httpd.serve_forever()
except KeyboardInterrupt:
    pass

httpd.server_close()

#END
