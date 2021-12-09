// +build !static

package gfx

//#cgo windows LDFLAGS: -lSDL2 -lSDL2_gfx
//#cgo linux freebsd darwin openbsd pkg-config: sdl2
//#cgo linux freebsd darwin openbsd LDFLAGS: -lSDL2_gfx
import "C"
