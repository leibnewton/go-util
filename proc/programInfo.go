package proc

import (
    "os"
    "path/filepath"
)

const (
    UnknownName = "=unknown="
)

var (
    programPath string
    programName string
)

func init() {
    ex, err := os.Executable()
    if err != nil {
        programPath = "."
        programName = UnknownName
        return
    }
    ultimPath, err := filepath.EvalSymlinks(ex)
    if err != nil {
        ultimPath = ex
    }
    programPath = filepath.Dir(ultimPath)
    programName = filepath.Base(ultimPath)
}

func ProgramPath() string {return programPath}
func ProgramName() string {return programName}

func ToAbsPath(relpath string) string {
    if filepath.IsAbs(relpath) {
        return relpath
    }
    return filepath.Join(programPath, relpath)
}
