#!/usr/bin/env python2.7

import sys
from PyQt4.QtCore import *
from PyQt4.QtGui import *
from PyQt4.QtWebKit import *

app = QApplication(sys.argv)

class htmlViewer(QWebView):
    def __init__(self,url, parent=None):
        QWebView.__init__(self,parent)
        self.setZoomFactor(1)
        self.setUrl(QUrl(url))
        self.printer = QPrinter(QPrinterInfo.defaultPrinter(),QPrinter.HighResolution)
        self.printer.setOutputFormat(QPrinter.PdfFormat)
        self.printer.setOrientation(QPrinter.Portrait)
        self.printer.setPaperSize(QPrinter.A4)
        self.printer.setFullPage(True)
        #self.printer.setResolution(72)
        self.printer.setOutputFileName("wkhtml2pdf-test.pdf")
        self.loadFinished.connect(self.execpreview)

    def execpreview(self,arg):
        self.print_(self.printer)

a = htmlViewer("https://www.soft112.com/")
a.show()

sys.exit(app.exec_())

### Creating Page Breaks
## @media print { #manualContent h2 {page-break-before:always;} }
