
// GO Lang

package main

import (
	"os"

	"github.com/ajstarks/svgo"
	"github.com/aaronarduino/goqrsvg"
	"github.com/boombuler/barcode/qr"
)

func main() {

	s := svg.New(os.Stdout)

	// Create the barcode
	qrCode, _ := qr.Encode("Hello World", qr.M, qr.Auto)

	// Write QR code to SVG
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	s.End()

}

// #END
