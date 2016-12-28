package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (socket *RedisSocket) LoadAll() {
	conn := socket.GetConn()
	defer func() {
		conn.Close()
		log.Println("end proccess charging_pile")
	}()
	conn.Do("SELECT", 1)
	var value interface{}
	var cursor_keys []interface{}
	var cursor string = "0"
	var keys []string
	var e error
	for {
		value, e = conn.Do("SCAN", cursor)
		CheckError(e)
		cursor_keys, e = redis.Values(value, e)
		CheckError(e)
		cursor, e = redis.String(cursor_keys[0], nil)
		CheckError(e)
		keys, e = redis.Strings(cursor_keys[1], nil)
		CheckError(e)
		socket.PipelineGetValue(keys)
		if cursor == "0" {
			return
		}
	}
}

func (socket *RedisSocket) PipelineGetValue(keys []string) {
	conn := socket.GetConn()
	defer func() {
		conn.Close()
	}()
	conn.Do("SELECT", 1)

	var index int = 0
	var key string = ""
	for index, key = range keys {
		conn.Send("GET", key)
		log.Println(key)
	}

	conn.Flush()

	for i := 0; i < index+1; i++ {
		v_redis, err := conn.Receive()

		if err != nil {
			log.Println(err)
			continue
		}

		v, _ := redis.Bytes(v_redis, nil)

		redis_pile_status := &Report.ChargingPileStatus{}
		err = proto.Unmarshal(v, redis_pile_status)
		if err != nil {
			log.Println("unmarshal error PipelineGetValue")
		} else {
			if redis_pile_status.Status != Report.ChargingPileStatus_MAINTAINACE {
				GetStatusChecker().Insert(redis_pile_status.Cpid, redis_pile_status.Timestamp, time.Now().Unix(), redis_pile_status.Id, redis_pile_status.StationId)
			}
		}
	}
	conn.Do("")
}
