// +build !static

package mix

//#cgo windows LDFLAGS: -lSDL2 -lSDL2_mixer
//#cgo linux freebsd darwin openbsd pkg-config: sdl2
//#cgo linux freebsd darwin openbsd LDFLAGS: -lSDL2_mixer
import "C"
