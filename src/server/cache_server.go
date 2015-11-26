// main
package main

import (
	"cache"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	mylog = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	addr := ""
	cs := NewTcpServer(addr)
	//default expiration time is 5 minites, cleanup time 30 seconds
	c := cache.New(5*time.Minute, 30*time.Second)
	if ok, _ := cs.Start(c); ok != 0 {
		log.Println("go-cache started failed")
		os.Exit(1)
	}
	log.Println("go-cache started and listening on 9891")

	<-signalChan
	fmt.Println("receive signalChan")
	cs.Stop()
	os.Exit(0)
}
