// +build !windows

package notify

func ShowMessage(bt BoxType, title, message string) {}

func ShowSysTopMessage(bt BoxType, title, message string) {}

func ShowAppTopMessage(bt BoxType, title, message string) {}

func ShowConfirm(bt BoxType, title, message string) bool {
	return true
}

func ShowSelection(bt BoxType, title, message string) bool {
	return true
}
