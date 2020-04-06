package dump

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/leibnewton/go-util/proc"
)

var (
	dumpPath string = "."
	maxDumps int    = 10
)

func SetPath(max int, dpath string) error {
	dumpPath = proc.ToAbsPath(dpath)
	if err := os.MkdirAll(dumpPath, 0777); err != nil {
		return err
	}
	if max < 1 {
		max = 1
	}
	maxDumps = max
	return nil
}

func getFileName(idx int) string {
	lastPart := "meta"
	if idx >= 0 {
		lastPart = strconv.Itoa(idx)
	}
	fname := fmt.Sprintf("dump-%s-%s.log", proc.ProgramName(), lastPart)
	return filepath.Join(dumpPath, fname)
}

type dumpMeta struct {
	Index int `json:"index"`
}

func getDumpName() (string, error) {
	fname := getFileName(-1)
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", errors.Wrap(err, "open failed")
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "read failed")
	}
	var meta dumpMeta
	if len(content) > 0 {
		if err = yaml.Unmarshal(content, &meta); err != nil {
			log.Printf("meta corrupted: %v", err)
		} else {
			meta.Index++
		}
		if meta.Index >= maxDumps {
			meta.Index = 0
		}
	}
	content, err = yaml.Marshal(&meta)
	if err != nil {
		return "", errors.Wrap(err, "marshal meta failed")
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", errors.Wrap(err, "seek to start failed")
	}
	n, err := f.Write(content)
	if err != nil {
		return "", errors.Wrap(err, "write failed")
	}
	err = f.Truncate(int64(n))
	return getFileName(meta.Index), errors.Wrap(err, "truncate failed")
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
