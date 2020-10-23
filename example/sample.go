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

	fmt.Println("hello world:", os.Args)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "panic":
			defer dump.PanicHandler()
			doTask()
		case "gopanic":
			go dump.WithPanicHandler(doTask)
			time.Sleep(time.Second)
		case "recover":
			defer dump.RecoverHandler()
			doTask()
		case "gorecover":
			go dump.WithRecoverHandler(doTask)
			time.Sleep(time.Second)
		}
		select {}
	}
	log.Printf("done")
}
