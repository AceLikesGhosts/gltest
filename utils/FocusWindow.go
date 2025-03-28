package utils

import (
	"fmt"
	"syscall"
)

func FocusWindow(hwnd syscall.Handle) error {
	// check if its minimised, if so unminimise
	isIconic, _, _ := procIsIconic.Call(uintptr(hwnd))
	if isIconic != 0 {
		ret, _, err := procShowWindow.Call(uintptr(hwnd), uintptr(SW_RESTORE))
		if ret == 0 {
			return fmt.Errorf("ShowWindow failed: %v", err)
		}
	}

	// bring to front
	ret, _, err := procSetForegroundWindow.Call(uintptr(hwnd))
	if ret == 0 {
		return fmt.Errorf("SetForegroundWindow failed: %v", err)
	}

	procSetActiveWindow.Call(uintptr(hwnd))
	return nil
}
