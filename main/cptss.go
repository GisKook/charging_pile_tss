package main

import (
	"fmt"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/mq"
	"github.com/giskook/charging_pile_tss/redis_socket"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := conf.ReadConfig("./conf.json")

	checkError(err)
	// create a redis  socket
	redis_socket, e := redis_socket.NewRedisSocket(configuration.Redis)
	checkError(e)
	// create a mq socket
	mq_socket := mq.GetNsqSocket(configuration.Nsq)
	mq_socket.Start()
	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
	redis_socket.Close()
	mq_socket.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
