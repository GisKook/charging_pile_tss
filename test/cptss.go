package main

import (
	"fmt"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/giskook/charging_pile_tss/redis_socket"
	"github.com/giskook/charging_pile_tss/server"
	"github.com/golang/protobuf/proto"
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
	reply, e := conn.Do("select", 1)
	log.Println(reply)
	log.Println(e)

	pile := &Report.ChargingPileStatus{
		DasUuid:   "das",
		Cpid:      1000000000000050,
		Id:        50,
		StationId: 14,
	}
	data, _ := proto.Marshal(pile)
	log.Println(data)

	conn.Do("SET", "14.50.1000000000000050", data)
	//	value, _ := conn.Do("GET", 1024)
	//	v, _ := redis.Bytes(value, nil)
	//
	//	redis_cps := &Report.ChargingStationStatus{}
	//	proto.Unmarshal(v, redis_cps)
	//	log.Println(redis_cps)

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
