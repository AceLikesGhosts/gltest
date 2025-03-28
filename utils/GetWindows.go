package utils

import (
	"syscall"
	"unsafe"
)

type WindowInfo struct {
	Hwnd  syscall.Handle
	Title string
}

// for GetWindowRect
type RECT struct {
	Left, Top, Right, Bottom int32
}

func enumWindowsProc(hwnd syscall.Handle, lparam uintptr) uintptr {
	// language server isnt very happy about this
	//  but this is completely valid
	windowInfoList := (*[]WindowInfo)(unsafe.Pointer(lparam))

	// skip non-visible windows
	ret, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
	if ret == 0 {
		return 1
	}

	// skip if has owner window (not top-level unowned)
	owner, _, _ := procGetWindow.Call(uintptr(hwnd), GW_OWNER)
	if owner != 0 {
		return 1
	}

	exStyle, _, _ := procGetWindowLongPtrW.Call(uintptr(hwnd), uintptr(GWL_EXSTYLE))

	// skip tool windows
	if (exStyle & WS_EX_TOOLWINDOW) != 0 {
		return 1
	}

	// get window title
	buf := make([]uint16, 256)
	procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	title := syscall.UTF16ToString(buf)

	if title == "" {
		return 1
	}

	style, _, _ := procGetWindowLongPtrW.Call(uintptr(hwnd), uintptr(GWL_STYLE))
	hasCaption := (style & WS_CAPTION) != 0
	hasThickFrame := (style & WS_THICKFRAME) != 0
	isAppWindow := (exStyle & WS_EX_APPWINDOW) != 0

	// include window if it's either:
	// 1. marked as an app window, OR
	// 2. has proper window decorations (caption and frame)
	if isAppWindow || (hasCaption && hasThickFrame) {
		*windowInfoList = append(*windowInfoList, WindowInfo{
			Hwnd:  hwnd,
			Title: title,
		})
	}

	return 1
}

func GetWindowTitles() ([]WindowInfo, error) {
	windowInfoList := []WindowInfo{}

	enumWindowsProcPtr := syscall.NewCallback(enumWindowsProc)
	procEnumWindows.Call(enumWindowsProcPtr, uintptr(unsafe.Pointer(&windowInfoList)))
	// return nil, fmt.Errorf("test")

	return windowInfoList, nil
}
