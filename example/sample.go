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
	err := dump.SetPath(3, "dmp")
	if err != nil {
		log.Fatal(err)
	}
	defer dump.PanicHandler()

	fmt.Println("hello world")
	go dump.WithPanicHandler(doTask)
	time.Sleep(time.Second)
}
