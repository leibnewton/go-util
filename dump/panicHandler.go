package dump

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/leibnewton/go-util/notify"
	"github.com/leibnewton/go-util/proc"
)

var (
	dumpPath    = "."
	maxDumps    = 10
	showMsgBox  = true
	titleMsgBox = "程序出现异常(Exception Caught)"

	panicReporter atomic.Value
)

type PanicReporter func(interface{}, []byte) // error, stack

// if relativeToWorkDir set to true, will save dump file according to working directory,
//    otherwise will save dump file according to directory where the exe resides.
func SetPath(max int, dpath string, relativeToWorkDir bool) error {
	if relativeToWorkDir {
		dumpPath, _ = filepath.Abs(dpath)
	} else {
		dumpPath = proc.ToAbsPath(dpath)
	}
	log.Printf("dumps will be saved to %s", dumpPath)
	if err := os.MkdirAll(dumpPath, 0777); err != nil {
		return err
	}
	if max < 1 {
		max = 1
	}
	maxDumps = max
	return nil
}

func EnableMessageBox(enable bool, title string) {
	showMsgBox = enable
	if len(title) > 0 {
		titleMsgBox = title
	}
}

// set log output writer
func SetLogOutput(logWriter io.Writer) {
	log.SetOutput(logWriter)
}

func SetPanicReporter(reporter PanicReporter) {
	panicReporter.Store(reporter)
}

// dump then panic
func PanicHandler() {
	if err := recover(); err != nil { // NOTE: cannot put recover inside `panicHandler`, otherwise no effect.
		panicHandler(err, true)
	}
}

// dump then recover
func RecoverHandler() {
	if err := recover(); err != nil { // NOTE: cannot put recover inside `panicHandler`, otherwise no effect.
		panicHandler(err, false)
	}
}

func panicHandler(err interface{}, passPanic bool) {
	defer func() {
		if passPanic {
			if showMsgBox {
				notify.ShowSysTopMessage(notify.BoxTypeError, titleMsgBox,
					fmt.Sprintf("Application will EXIT due to error:\n  %v", err))
			}
			panic(err)
		} else {
			log.Printf("PanicHandler: dump panic and continue. detail: %v", err)
		}
	}()

	fname, ierr := getDumpName()
	if ierr != nil {
		log.Printf("PanicHandler: getDumpName failed: %v", ierr)
		return
	}
	log.Printf("dump to file %s", fname)

	f, ierr := os.Create(fname)
	if ierr != nil {
		log.Printf("PanicHandler: create %s failed: %v", fname, ierr)
		return
	}
	defer f.Close()
	header := fmt.Sprintf(`Time: %s
Pid: %d
Reason: %+v
===================
`, time.Now().Format("2006-01-02 15:04:05.000 MST"), os.Getpid(), err)
	stack := debug.Stack()
	f.WriteString(header) //输出panic信息
	f.Write(stack)        //输出堆栈信息
	reporter := panicReporter.Load()
	if rep, ok := reporter.(PanicReporter); ok && rep != nil {
		rep(err, stack)
	}
	return
}

// dump then panic goroutine
func WithPanicHandler(routine func()) {
	defer PanicHandler()
	routine()
}

// dump then recover goroutine
func WithRecoverHandler(routine func()) {
	defer RecoverHandler()
	routine()
}
