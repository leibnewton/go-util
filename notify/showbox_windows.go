package notify

import (
	"syscall"
	"unsafe"
)

var (
	user32Dll  = syscall.NewLazyDLL("user32.dll")
	messageBox = user32Dll.NewProc("MessageBoxW")
)

const (
	MB_OK                uint = 0x00000000
	MB_OKCANCEL          uint = 0x00000001
	MB_ABORTRETRYIGNORE  uint = 0x00000002
	MB_YESNOCANCEL       uint = 0x00000003
	MB_YESNO             uint = 0x00000004
	MB_RETRYCANCEL       uint = 0x00000005
	MB_CANCELTRYCONTINUE uint = 0x00000006

	MB_ICONHAND        uint = 0x00000010
	MB_ICONERROR       uint = MB_ICONHAND
	MB_ICONQUESTION    uint = 0x00000020
	MB_ICONEXCLAMATION uint = 0x00000030
	MB_ICONWARNING     uint = MB_ICONEXCLAMATION
	MB_ICONASTERISK    uint = 0x00000040
	MB_ICONINFORMATION uint = MB_ICONASTERISK

	MB_SYSTEMMODAL uint = 0x00001000
	MB_TOPMOST     uint = 0x00040000

	IDOK       = 1  // The OK button was selected.
	IDCANCEL   = 2  // The Cancel button was selected.
	IDABORT    = 3  // The Abort button was selected.
	IDRETRY    = 4  // The Retry button was selected.
	IDIGNORE   = 5  // The Ignore button was selected.
	IDYES      = 6  // The Yes button was selected.
	IDNO       = 7  // The No button was selected.
	IDTRYAGAIN = 10 // The Try Again button was selected.
	IDCONTINUE = 11 // The Continue button was selected.
)

func joinFlags(flags ...uint) uint {
	for _, flag := range flags {
		flags[0] |= flag
	}
	return flags[0]
}

func flagFromBoxType(bt BoxType) uint {
	switch bt {
	case BoxTypeWarn:
		return MB_ICONWARNING
	case BoxTypeError:
		return MB_ICONERROR
	case BoxTypeInfo:
		return MB_ICONINFORMATION
	}
	return MB_ICONINFORMATION
}

func MessageBox(hwnd uintptr, text, title string, flags uint) int {
	ret, _, _ := messageBox.Call(
		hwnd,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))
	return int(ret)
}

func ShowMessage(bt BoxType, title, message string) {
	flags := joinFlags(MB_OK, flagFromBoxType(bt))
	MessageBox(0, message, title, flags)
}

func ShowAppTopMessage(bt BoxType, title, message string) {
	flags := joinFlags(MB_OK, MB_TOPMOST, flagFromBoxType(bt))
	MessageBox(0, message, title, flags)
}

func ShowSysTopMessage(bt BoxType, title, message string) {
	flags := joinFlags(MB_OK, MB_SYSTEMMODAL, flagFromBoxType(bt))
	MessageBox(0, message, title, flags)
}

func ShowConfirm(bt BoxType, title, message string) bool {
	flags := joinFlags(MB_OKCANCEL, flagFromBoxType(bt))
	return MessageBox(0, message, title, flags) == IDOK
}

func ShowSelection(bt BoxType, title, message string) bool {
	flags := joinFlags(MB_YESNO, flagFromBoxType(bt))
	return MessageBox(0, message, title, flags) == IDYES
}
