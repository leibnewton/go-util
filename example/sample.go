package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/leibnewton/go-util/dump"
)

func doTask() {
	log.Printf("now is %v", time.Now())
	panic("show me stack trace")
}

func main() {
	err := dump.SetPath(3, "dmp", true)
	if err != nil {
		log.Fatal(err)
	}
	defer dump.PanicHandler()

	fmt.Println("hello world:", os.Args)
	if len(os.Args) > 1 {
		if os.Args[1] == "panic" {
			doTask()
		} else if os.Args[1] == "gopanic" {
			go dump.WithPanicHandler(doTask)
			time.Sleep(time.Second)
		}
	}
}
