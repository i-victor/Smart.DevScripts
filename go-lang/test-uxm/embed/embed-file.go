
// GO Lang

package main

import (
	"embed"
)

func main() {

	//go:embed hello.txt
	var f embed.FS
	data, _ := f.ReadFile("hello.txt")
	print(string(data))

}

// #END
