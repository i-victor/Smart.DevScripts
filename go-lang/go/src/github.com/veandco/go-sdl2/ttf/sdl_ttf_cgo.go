// +build !static

package ttf

//#cgo windows LDFLAGS: -lSDL2 -lSDL2_ttf
//#cgo linux freebsd darwin openbsd pkg-config: SDL2_ttf
//#cgo linux freebsd darwin openbsd LDFLAGS: -lSDL2_ttf
import "C"
