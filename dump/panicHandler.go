package dump

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/leibnewton/go-util/proc"
)

var (
	dumpPath string = "."
	maxDumps int    = 10
)

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

func PanicHandler() {
	if err := recover(); err != nil {
		fname, ierr := getDumpName()
		if ierr != nil {
			log.Printf("PanicHandler: getDumpName failed: %v", ierr)
			panic(err)
			return
		}
		log.Printf("dump to file %s", fname)

		f, ierr := os.Create(fname)
		if ierr != nil {
			log.Printf("PanicHandler: create %s failed: %v", fname, ierr)
			panic(err)
			return
		}
		defer f.Close()
		header := fmt.Sprintf(`Time: %s
Pid: %d
Reason: %+v
===================
`, time.Now().Format("2006-01-02 15:04:05.000 MST"), os.Getpid(), err)
		f.WriteString(header)  //输出panic信息
		f.Write(debug.Stack()) //输出堆栈信息
		panic(err)
	}
}

func WithPanicHandler(routine func()) {
	defer PanicHandler()
	routine()
}
