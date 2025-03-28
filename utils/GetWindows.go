package utils

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

type WindowInfo struct {
	Hwnd  syscall.Handle
	Title string
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
		var pid uint32

		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
		hProcess, _, _ := procOpenProcess.Call(PROCESS_QUERY_LIMITED_INFORMATION|PROCESS_VM_READ, 0, uintptr(pid))

		if hProcess != 0 {
			defer procCloseHandle.Call(hProcess)

			exeBuf := make([]uint16, 260)
			procGetModuleFileNameExW.Call(hProcess, 0, uintptr(unsafe.Pointer(&exeBuf[0])), uintptr(len(exeBuf)))
			exePath := syscall.UTF16ToString(exeBuf)

			if exePath != "" {
				exeName := filepath.Base(exePath)
				// Combine executable name with window title
				title = exeName + " - " + title
			}
		}

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

	return windowInfoList, nil
}
