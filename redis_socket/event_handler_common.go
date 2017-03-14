package redis_socket

import (
	//"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/charging_pile_tss/base"
	"github.com/giskook/charging_pile_tss/db"
	"github.com/giskook/charging_pile_tss/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"strconv"
)

const (
	PROTOCOL_CHARGE_PILE_STATUS_IDLE     uint8 = 0
	PROTOCOL_CHARGE_PILE_STATUS_CHARGING uint8 = 1
	PROTOCOL_CHARGE_PILE_STATUS_STARTED  uint8 = 100
	PROTOCOL_CHARGE_PILE_STATUS_STOPPED  uint8 = 101
	PROTOCOL_CHARGE_PILE_STATUS_OFFLINE  uint8 = 255

	PROTOCOL_CHARGE_PILE_STATUS_COMPELETED  uint8 = 6 // for db
	PROTOCOL_CHARGE_PILE_STATUS_CHARGING_DB uint8 = 5
)

type charge_pile_status struct {
	old_status uint8
	new_status uint8
	status     *Report.ChargingPileStatus
}

func GetKey(command *Report.ChargingPileStatus) string {
	station_id := strconv.FormatUint(uint64(command.StationId), 10)
	id := strconv.FormatUint(uint64(command.Id), 10)
	cpid := strconv.FormatUint(command.Cpid, 10)
	return station_id + "." + id + "." + cpid
}

func (socket *RedisSocket) ProcessChargingPile() {
	socket.Mutex_ChargingPiles.Lock()
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

		transactions := make(chan *base.TransactionDetail, 1024)

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
					old_status := redis_pile.Status
					socket.ProccessIncomingStatus(transactions, redis_pile, socket.ChargingPiles[i])
					log.Println(redis_pile)

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
		db.GetDBSocket().TransactionChan <- transactions

		close(transactions)

		socket.ChargingPiles = socket.ChargingPiles[:0]

		for _, new_pkg := range tobe_commit_cps {
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

func (socket *RedisSocket) ProccessIncomingStatus(ch chan *base.TransactionDetail, redis_pile *Report.ChargingPileStatus, new_status *Report.ChargingPileStatus) {
	redis_pile.Timestamp = new_status.Timestamp
	redis_pile.DasUuid = new_status.DasUuid

	redis_pile.ChargingDuration = new_status.ChargingDuration
	redis_pile.ChargingCapacity = new_status.ChargingCapacity
	redis_pile.Status = new_status.Status
	if new_status.Status == uint32(PROTOCOL_CHARGE_PILE_STATUS_STARTED) {
		log.Println("charging started")
		log.Println(new_status.StartTime)
		redis_pile.StartTime = new_status.StartTime
		redis_pile.ChargingCost = 0.0
		redis_pile.ChargingCostE = 0.0
		redis_pile.StartMeterReading = new_status.StartMeterReading
		redis_pile.EndMeterReading = new_status.StartMeterReading
		redis_pile.CurrentOrderNumber = new_status.CurrentOrderNumber
		redis_pile.PreCharge = new_status.PreCharge

		ch <- &base.TransactionDetail{
			TransactionID:     redis_pile.CurrentOrderNumber,
			Status:            PROTOCOL_CHARGE_PILE_STATUS_CHARGING_DB,
			ChargingDuration:  0,
			ChargingCapacity:  0,
			ChargingCost:      0,
			ChargingCostEle:   0,
			StartTime:         redis_pile.StartTime,
			StartMeterReading: redis_pile.StartMeterReading,
		}
	} else if new_status.Status == uint32(PROTOCOL_CHARGE_PILE_STATUS_CHARGING) {
		redis_pile.RealTimeCurrent = new_status.RealTimeCurrent
		redis_pile.RealTimeVoltage = new_status.RealTimeVoltage
		redis_pile.CurrentOrderNumber = new_status.CurrentOrderNumber
		redis_pile.ChargingCost, redis_pile.ChargingCostE = CalcCost(redis_pile.StationId, redis_pile.ChargingCost, redis_pile.ChargingCostE, redis_pile.EndMeterReading, new_status.EndMeterReading, new_status.Timestamp)
		redis_pile.EndMeterReading = new_status.EndMeterReading
		socket.ChargingCost <- &base.ChargingCost{
			Uuid: redis_pile.DasUuid,
			Tid:  redis_pile.Cpid,
			Cost: uint32(redis_pile.ChargingCost * 100),
		}
		if redis_pile.ChargingCost >= redis_pile.PreCharge {
			socket.StopChargingNotifyChan <- &base.StopChargingNotify{
				Uuid: redis_pile.DasUuid,
				Tid:  redis_pile.Cpid,
			}
		}
	} else if new_status.Status == uint32(PROTOCOL_CHARGE_PILE_STATUS_STOPPED) {
		redis_pile.ChargingCost, redis_pile.ChargingCostE = CalcCost(redis_pile.StationId, redis_pile.ChargingCost, redis_pile.ChargingCostE, redis_pile.EndMeterReading, new_status.EndMeterReading, new_status.Timestamp)
		redis_pile.EndMeterReading = new_status.EndMeterReading
		redis_pile.EndTime = new_status.EndTime
		redis_pile.RealTimeCurrent = 0.0
		redis_pile.RealTimeVoltage = 0.0
		log.Println("Charge stopped")
		log.Println(redis_pile.StartTime)
		log.Println(redis_pile.EndTime)

		ch <- &base.TransactionDetail{
			TransactionID:     redis_pile.CurrentOrderNumber,
			Status:            PROTOCOL_CHARGE_PILE_STATUS_COMPELETED,
			ChargingDuration:  redis_pile.ChargingDuration,
			ChargingCapacity:  redis_pile.ChargingCapacity,
			ChargingCost:      redis_pile.ChargingCost,
			ChargingCostEle:   redis_pile.ChargingCostE,
			StartTime:         redis_pile.StartTime,
			EndTime:           redis_pile.EndTime,
			StartMeterReading: redis_pile.StartMeterReading,
			EndMeterReading:   redis_pile.EndMeterReading,
		}
	}

}
