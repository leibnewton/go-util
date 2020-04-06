package dump

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/leibnewton/go-util/proc"
)

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
