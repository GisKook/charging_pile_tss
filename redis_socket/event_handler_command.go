package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"strconv"
)

type charge_pile_status struct {
	old_status uint8
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

func (socket *RedisSocket) UpdateChargeStation(pile_status []*charge_pile_status) {
	conn := socket.GetConn()
	defer conn.Close()
	conn.Do("SELECT", 0)
	map_station_status := make(map[uint32]*charge_station_status)

	for _, pile := range pile_status {
		station_id := pile.status.StationId
		_, ok := map_station_status[station_id]
		if !ok {
			map_station_status[station_id] = &charge_station_status{
				idle_pile:     0,
				charging_pile: 0,
				error_pile:    0,
			}
		}
		if pile.old_status == 0 {
			map_station_status[station_id].idle_pile--
		} else if pile.old_status == 1 ||
			pile.old_status == 2 ||
			pile.old_status == 3 {
			map_station_status[station_id].charging_pile--
		} else if pile.old_status == 4 ||
			pile.old_status == 5 {
			map_station_status[station_id].error_pile--
		}

		if pile.old_status == 0 {
			map_station_status[station_id].idle_pile++
		} else if pile.old_status == 1 ||
			pile.old_status == 2 ||
			pile.old_status == 3 {
			map_station_status[station_id].charging_pile++
		} else if pile.old_status == 4 ||
			pile.old_status == 5 {
			map_station_status[station_id].error_pile++
		}
	}

	for index, _ := range map_station_status {
		conn.Send("GET", index)
	}

	conn.Flush()

	tobe_commit_stations := make([]*Report.ChargingStationStatus, 0)
	for _, _station := range map_station_status {
		v_redis, err := conn.Receive()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		v, _ := redis.Bytes(v_redis, nil)

		redis_stations := &Report.ChargingStationStatus{}
		err = proto.Unmarshal(v, redis_stations)
		if err != nil {
			log.Println("unmarshal error")
		} else {
			redis_stations.FreePileNumber += _station.idle_pile
			redis_stations.WorkingPileNumber += _station.charging_pile
			redis_stations.ErrorPileNumber += _station.error_pile
			tobe_commit_stations = append(tobe_commit_stations,
				redis_stations)
		}
	}
	for _, new_pkg := range tobe_commit_stations {
		data, _ := proto.Marshal(new_pkg)
		conn.Send("SET", new_pkg.Id, data)
		new_pkg = nil
	}
	conn.Flush()
	conn.Do("EXEC")
	conn.Close()
}
func (socket *RedisSocket) ProcessChargingPile() {
	log.Println("proccess charging pile")
	conn := socket.GetConn()
	defer conn.Close()
	conn.Do("SELECT", 1)

	var index int = 0
	var pkg *Report.ChargingPileStatus
	log.Println(socket.ChargingPiles)
	for index, pkg = range socket.ChargingPiles {
		conn.Send("GET", GetKey(pkg))
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
			log.Println("unmarshal error")
		} else {
			if redis_pile.Timestamp < socket.ChargingPiles[i].Timestamp {
				redis_pile.Timestamp = socket.ChargingPiles[i].Timestamp
				redis_pile.DasUuid = socket.ChargingPiles[i].DasUuid

				old_status := redis_pile.Status
				redis_pile.Status = socket.ChargingPiles[i].Status
				redis_pile.ChargingDuration = socket.ChargingPiles[i].ChargingDuration
				redis_pile.ChargingCapacity = socket.ChargingPiles[i].ChargingCapacity
				redis_pile.ChargingPrice = socket.ChargingPiles[i].ChargingPrice

				tobe_commit_cps = append(tobe_commit_cps,

					&charge_pile_status{
						old_status: uint8(old_status),
						status:     redis_pile})
			}
			socket.ChargingPiles[i] = nil
		}
	}

	socket.ChargingPiles = socket.ChargingPiles[:0]

	for _, new_pkg := range tobe_commit_cps {
		data, _ := proto.Marshal(new_pkg.status)
		conn.Send("SET", GetKey(new_pkg.status), data)
		log.Println(GetKey(new_pkg.status))
		new_pkg = nil
	}
	conn.Flush()
	conn.Do("EXEC")
	conn.Close()
	socket.UpdateChargeStation(tobe_commit_cps)

	tobe_commit_cps = tobe_commit_cps[:0]
}
