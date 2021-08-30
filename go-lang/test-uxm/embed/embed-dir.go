
// GO Lang

package main

import (
	"embed"
)

// https://pkg.go.dev/embed

func main() {

	// content holds our static web server content.
	//go:embed image/* template/*
	//go:embed html/index.html
	var content embed.FS

}

// #END
