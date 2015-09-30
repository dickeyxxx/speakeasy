// +build windows

package speakeasy

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func getPassword() (string, error) {
	pass := []byte{}
	for {
		if v := getch(); v == 127 || v == 8 {
			if l := len(pass); l > 0 {
				pass = pass[:l-1]
			}
		} else if v == 13 || v == 10 {
			break
		} else if v != 0 {
			pass = append(pass, v)
		}
	}
	println()
	return string(pass), nil
}

func getch() byte {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procReadConsole := modkernel32.NewProc("ReadConsoleW")
	procGetConsoleMode := modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := modkernel32.NewProc("SetConsoleMode")

	var mode uint32
	pMode := &mode
	procGetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pMode)))

	var echoMode, lineMode uint32
	echoMode = 4
	lineMode = 2
	var newMode uint32
	newMode = mode ^ (echoMode | lineMode)

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(newMode))

	line := make([]uint16, 1)
	pLine := &line[0]
	var n uint16
	procReadConsole.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pLine)), uintptr(len(line)), uintptr(unsafe.Pointer(&n)))

	b := []byte(string(utf16.Decode(line)))

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(mode))

	// Not sure how this could happen, but it did for someone
	if len(b) > 0 {
		return b[0]
	}
	return 13
}
