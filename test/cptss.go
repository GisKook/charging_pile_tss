package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
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
	"time"
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

	station := &Report.ChargingStationStatus{
		Timestamp:              uint64(time.Now().Unix()),
		Id:                     1024,
		Name:                   "cetcnav",
		City:                   130100,
		Province:               130000,
		Area:                   130123,
		Address:                "luquan",
		Contacter:              "james",
		Mobile:                 "13732143001",
		Lat:                    38.1,
		Lng:                    108.1,
		OperatorId:             1,
		OperatorScale:          3,
		PileNumber:             16,
		FreePileNumber:         4,
		WorkingPileNumber:      10,
		ErrorPileNumber:        2,
		ChargeMoneyToday:       1000.3,
		ChargeElectricityToday: 500.8,
	}
	data, _ := proto.Marshal(station)
	log.Println(data)

	conn.Do("SET", 1024, data)
	value, _ := conn.Do("GET", 1024)
	v, _ := redis.Bytes(value, nil)

	redis_cps := &Report.ChargingStationStatus{}
	proto.Unmarshal(v, redis_cps)
	log.Println(redis_cps)

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
