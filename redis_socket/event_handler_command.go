package redis_socket

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"strconv"
)

type charge_pile_status struct {
	old_status uint8
	new_status uint8
	status     *Report.ChargingPileStatus
}

type charge_station_status struct {
	idle_pile     int32
	charging_pile int32
	error_pile    int32
}

func GetKey(command *Report.ChargingPileStatus) string {
	station_id := strconv.FormatUint(uint64(command.StationId), 10)
	id := strconv.FormatUint(uint64(command.Id), 10)
	cpid := strconv.FormatUint(command.Cpid, 10)
	return station_id + "." + id + "." + cpid
}

func (socket *RedisSocket) GetStationID(pile_status []*charge_pile_status) []uint32 {

	var station_ids []uint32
	for _, p := range pile_status {
		if p.old_status == p.new_status {
			continue
		}
		for _, station_id := range station_ids {
			if station_id != p.status.StationId {
				station_ids = append(station_ids, p.status.StationId)
			}
		}
	}

	return station_ids
}

func (socket *RedisSocket) CalcSingleStation(station_id uint32) (int32, int32, int32, int32) {
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
			switch redis_pile_status.Status {
			case Report.ChargingPileStatus_MAINTAINACE:
				error_pile += 1
			case Report.ChargingPileStatus_OFFLINE:
				error_pile += 1
			case Report.ChargingPileStatus_IDLE:
				idle_pile += 1
			case Report.ChargingPileStatus_CHARGING:
				charging_pile += 1
			case Report.ChargingPileStatus_TOBECHARGING:
				charging_pile += 1
			case Report.ChargingPileStatus_FULL:
				charging_pile += 1
			}
		}
	}
	conn.Do("")

	return int32(pile_count), idle_pile, charging_pile, error_pile
}

func (socket *RedisSocket) UpdateSingleStation(station uint32, pile_count int32, idle_pile int32, charging_pile int32, error_pile int32) {
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
	station_ids := socket.GetStationID(pile_status)
	//	conn := socket.GetConn()
	//	defer conn.Close()
	//	conn.Do("SELECT", 0)

	for _, station_id := range station_ids {
		pile_count, idle_pile, charging_pile, error_pile := socket.CalcSingleStation(station_id)
		socket.UpdateSingleStation(station_id, pile_count, idle_pile, charging_pile, error_pile)

	}
}

//func (socket *RedisSocket) UpdateChargeStation(pile_status []*charge_pile_status) {
//	log.Println("update station")
//	log.Println(len(pile_status))
//	conn := socket.GetConn()
//	defer conn.Close()
//	conn.Do("SELECT", 0)
//	map_station_status := make(map[uint32]*charge_station_status)
//
//	for _, pile := range pile_status {
//		station_id := pile.status.StationId
//		_, ok := map_station_status[station_id]
//		if !ok {
//			map_station_status[station_id] = &charge_station_status{
//				idle_pile:     0,
//				charging_pile: 0,
//				error_pile:    0,
//			}
//		}
//		if pile.old_status != pile.new_status {
//			if pile.old_status == 0 {
//				map_station_status[station_id].idle_pile--
//			} else if pile.old_status == 1 {
//				map_station_status[station_id].charging_pile--
//			} else if pile.old_status == 4 {
//				map_station_status[station_id].error_pile--
//			}
//
//			if pile.new_status == 0 {
//				map_station_status[station_id].idle_pile++
//			} else if pile.new_status == 1 {
//				map_station_status[station_id].charging_pile++
//			} else if pile.new_status == 4 {
//				map_station_status[station_id].error_pile++
//			}
//		}
//	}
//
//	for index, _ := range map_station_status {
//		conn.Send("GET", index)
//		log.Printf("GET station %d\n", index)
//		log.Println(map_station_status[index])
//	}
//
//	conn.Flush()
//
//	tobe_commit_stations := make([]*Report.ChargingStationStatus, 0)
//	for _, _station := range map_station_status {
//		v_redis, err := conn.Receive()
//		if err != nil {
//			log.Println(err.Error())
//			continue
//		}
//		v, _ := redis.Bytes(v_redis, nil)
//
//		redis_stations := &Report.ChargingStationStatus{}
//		err = proto.Unmarshal(v, redis_stations)
//		if err != nil {
//			log.Println("update station unmarshal error")
//		} else {
//			//redis_stations.FreePileNumber = 2
//			//redis_stations.WorkingPileNumber = 0
//			//redis_stations.ErrorPileNumber = 0
//			//log.Println(_station)
//			redis_stations.FreePileNumber += _station.idle_pile
//			redis_stations.WorkingPileNumber += _station.charging_pile
//			redis_stations.ErrorPileNumber += _station.error_pile
//			tobe_commit_stations = append(tobe_commit_stations,
//				redis_stations)
//		}
//	}
//	for _, new_pkg := range tobe_commit_stations {
//		log.Println(new_pkg)
//		data, _ := proto.Marshal(new_pkg)
//		conn.Send("SET", new_pkg.Id, data)
//		new_pkg = nil
//	}
//	conn.Flush()
//	conn.Do("EXEC")
//}

func (socket *RedisSocket) ProcessChargingPile() {
	socket.Mutex_ChargingPiles.Lock()
	//log.Println("start proccess charging pile")
	conn := socket.GetConn()
	defer func() {
		conn.Close()
		socket.Mutex_ChargingPiles.Unlock()
	}()
	log.Println(len(socket.ChargingPiles))
	if len(socket.ChargingPiles) != 0 {
		conn.Do("SELECT", 1)

		var index int = 0
		var pkg *Report.ChargingPileStatus
		//log.Println(socket.ChargingPiles)
		for index, pkg = range socket.ChargingPiles {
			conn.Send("GET", GetKey(pkg))
			log.Println(GetKey(pkg))
		}

		conn.Flush()

		tobe_commit_cps := make([]*charge_pile_status, 0)
		for i := 0; i < index+1; i++ {
			v_redis, err := conn.Receive()

			if err != nil {
				log.Println(err.Error())
				continue
			}

			v, _ := redis.Bytes(v_redis, nil)

			redis_pile := &Report.ChargingPileStatus{}
			err = proto.Unmarshal(v, redis_pile)
			if err != nil {
				log.Println("ProcessChargingPile unmarshal error ")
			} else {
				//log.Println(redis_pile)
				if redis_pile.Timestamp <= socket.ChargingPiles[i].Timestamp {
					redis_pile.Timestamp = socket.ChargingPiles[i].Timestamp
					redis_pile.DasUuid = socket.ChargingPiles[i].DasUuid

					old_status := redis_pile.Status
					redis_pile.Status = socket.ChargingPiles[i].Status
					redis_pile.ChargingDuration = socket.ChargingPiles[i].ChargingDuration
					redis_pile.ChargingCapacity = socket.ChargingPiles[i].ChargingCapacity
					redis_pile.ChargingPrice = socket.ChargingPiles[i].ChargingPrice
					redis_pile.CurrentOrderNumber = socket.ChargingPiles[i].CurrentOrderNumber

					tobe_commit_cps = append(tobe_commit_cps,

						&charge_pile_status{
							old_status: uint8(old_status),
							new_status: uint8(redis_pile.Status),
							status:     redis_pile})
				}
				//log.Println(redis_pile)
				socket.ChargingPiles[i] = nil
			}
		}

		socket.ChargingPiles = socket.ChargingPiles[:0]

		for _, new_pkg := range tobe_commit_cps {
			//log.Println(new_pkg.status)
			data, _ := proto.Marshal(new_pkg.status)
			conn.Send("SET", GetKey(new_pkg.status), data)
			//log.Println(GetKey(new_pkg.status))
			new_pkg = nil
		}
		conn.Flush()
		conn.Do("EXEC")
		socket.UpdateChargeStation(tobe_commit_cps)

		tobe_commit_cps = tobe_commit_cps[:0]
	}
}
