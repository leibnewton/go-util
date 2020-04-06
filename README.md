# Golang utility

## Introduction

Provide common utilities in golang. Now supports panic handler, program info helpers.

## Installation

   To install, run:

   ```
$ go get github.com/leibnewton/go-util
   ```

   Example:

   ```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/leibnewton/go-util/dump"
)

func doTask() {
    log.Printf("now is %v", time.Now())
    panic("show me stack trace")
}

func main() {
  err := dump.SetPath(3, "dmp") // store in dmp directory, rotate in 3 files.
  if err != nil {
      log.Fatal(err)
  }
  defer dump.PanicHandler()

  fmt.Println("hello world")
  doTask()
}
   ```
