package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	test_chan := make(chan int, 10)
	go func() {
		for i := 0; i < 120; i++ {
			test_chan <- i
		}

	}()
	for {
		select {
		case p := <-test_chan:
			log.Println(p)
		}
	}

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}
