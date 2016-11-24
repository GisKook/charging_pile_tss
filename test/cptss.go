package main

import (
	"fmt"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/redis_socket"
	"github.com/giskook/charging_pile_tss/server"
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
	// create a mq socket
	mq_socket := server.NewNsqSocket(configuration.Nsq)
	// create a redis socket
	redis_socket, e := redis_socket.NewRedisSocket(configuration.Redis)
	checkError(e)
	// create server
	server := server.NewServer(mq_socket, redis_socket)
	server.Start()
	conn := redis_socket.GetConn()
	conn.Do("SET", "zhangkai", "world")

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
	server.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
