package utils

import "syscall"

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindow                = user32.NewProc("GetWindow")
	procGetWindowLongPtrW        = user32.NewProc("GetWindowLongPtrW")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procIsIconic                 = user32.NewProc("IsIconic")
	procShowWindow               = user32.NewProc("ShowWindow")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procSetActiveWindow          = user32.NewProc("SetActiveWindow")
	procDwmGetWindowAttribute    = user32.NewProc("DwmGetWindowAttribute")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess              = kernel32.NewProc("OpenProcess")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
	psapi                        = syscall.NewLazyDLL("psapi.dll")
	procGetModuleFileNameExW     = psapi.NewProc("GetModuleFileNameExW")
)

const (
	SW_RESTORE                        = 9 // restores the window to its previous position and size
	DWMWA_CLOAKED                     = 14
	WS_VISIBLE                        = 0x10000000
	WS_CAPTION                        = 0x00C00000
	WS_THICKFRAME                     = 0x00040000
	GW_OWNER                          = 4
	WS_EX_TOOLWINDOW                  = 0x00000080
	WS_EX_APPWINDOW                   = 0x00040000
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	PROCESS_VM_READ                   = 0x0010
)

var GWL_STYLE = -16
var GWL_EXSTYLE int32 = -20
