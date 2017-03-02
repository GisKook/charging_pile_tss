package redis_socket

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

func (socket *RedisSocket) GetStationID(pile_status []*charge_pile_status) []uint32 {
	if len(pile_status) == 0 {
		return nil
	}

	var station_ids []uint32
	for _, p := range pile_status {
		if p.old_status != p.new_status {
			log.Println("status not equal")
			log.Println(p.old_status)
			log.Println(p.new_status)
			if len(station_ids) == 0 {
				station_ids = append(station_ids, p.status.StationId)
				log.Println("add station")
			} else {
				log.Println("add station")
				for _, station_id := range station_ids {
					if station_id != p.status.StationId {
						station_ids = append(station_ids, p.status.StationId)
					}
				}
			}
		}
	}

	return station_ids
}

func (socket *RedisSocket) CalcSingleStation(station_id uint32) (int32, int32, int32, int32) {
	log.Println("CalcSingleStation")
	conn := socket.GetConn()
	defer conn.Close()
	conn.Do("SELECT", 1)

	var keys_format string = "%d.*.*"
	_keys := fmt.Sprintf(keys_format, station_id)
	keys_in_redis, _ := conn.Do("KEYS", _keys)
	keys, _ := redis.Strings(keys_in_redis, nil)

	var index int
	var key string
	for index, key = range keys {
		conn.Send("GET", key)
	}

	conn.Flush()

	pile_count := index + 1
	var idle_pile int32 = 0
	var charging_pile int32 = 0
	var error_pile int32 = 0
	for i := 0; i < pile_count; i++ {
		v_redis, err := conn.Receive()

		v, _ := redis.Bytes(v_redis, nil)
		redis_pile_status := &Report.ChargingPileStatus{}

		err = proto.Unmarshal(v, redis_pile_status)
		if err == nil {
			if redis_pile_status.Status == 0 {
				idle_pile += 1
			} else if redis_pile_status.Status == uint32(PROTOCOL_CHARGE_PILE_STATUS_OFFLINE) {
				error_pile += 1
			} else {
				charging_pile += 1
			}
		}
	}
	conn.Do("")

	return int32(pile_count), idle_pile, charging_pile, error_pile
}

func (socket *RedisSocket) UpdateSingleStation(station uint32, pile_count int32, idle_pile int32, charging_pile int32, error_pile int32) {
	log.Println("UpdateSingleStation")
	conn := socket.GetConn()
	defer conn.Close()
	conn.Do("SELECT", 0)
	station_bin, _ := conn.Do("GET", station)
	station_v, _ := redis.Bytes(station_bin, nil)

	redis_stations := &Report.ChargingStationStatus{}
	err := proto.Unmarshal(station_v, redis_stations)
	if err != nil {
		log.Println("update station unmarshal error")
	} else {
		redis_stations.FreePileNumber = idle_pile
		redis_stations.WorkingPileNumber = charging_pile
		redis_stations.ErrorPileNumber = error_pile
		redis_stations.PileNumber = pile_count
		data, _ := proto.Marshal(redis_stations)
		conn.Do("SET", station, data)
	}
}

func (socket *RedisSocket) UpdateChargeStation(pile_status []*charge_pile_status) {
	log.Println("UpdateChargeStation")
	station_ids := socket.GetStationID(pile_status)
	log.Println(station_ids)
	//	conn := socket.GetConn()
	//	defer conn.Close()
	//	conn.Do("SELECT", 0)

	if station_ids != nil {
		for _, station_id := range station_ids {
			pile_count, idle_pile, charging_pile, error_pile := socket.CalcSingleStation(station_id)
			log.Printf("station id %d total %d idle %d charging %d  error %d\n", station_id, pile_count, idle_pile, charging_pile, error_pile)
			socket.UpdateSingleStation(station_id, pile_count, idle_pile, charging_pile, error_pile)
		}
	}
}
