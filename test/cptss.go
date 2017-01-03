package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/giskook/charging_pile_tss/redis_socket"
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
	// create a redis socket
	redis_socket, e := redis_socket.NewRedisSocket(configuration.Redis)
	checkError(e)
	conn := redis_socket.GetConn()
	defer conn.Close()
	conn.Do("SELECT", 1)
	pile_bin, _ := conn.Do("GET", "24.84.1000000000000084")
	v, _ := redis.Bytes(pile_bin, nil)
	redis_pile := &Report.ChargingPileStatus{}
	err = proto.Unmarshal(v, redis_pile)
	if err != nil {
		log.Println(err)
	}
	log.Println(redis_pile)

	//redis_socket.LoadAll()
	//redis_socket.UpdateSingleStation(15)
	//	conn := redis_socket.GetConn()
	//	reply, e := conn.Do("select", 1)
	//	log.Println(reply)
	//	log.Println(e)
	//
	//	pile := &Report.ChargingPileStatus{
	//		DasUuid:   "das",
	//		Cpid:      1000000000000050,
	//		Id:        50,
	//		StationId: 14,
	//	}
	//	data, _ := proto.Marshal(pile)
	//	log.Println(data)
	//
	//	conn.Do("SET", "14.50.1000000000000050", data)

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
