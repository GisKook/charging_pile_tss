package redis_socket

import (
	"github.com/HuKeping/rbtree"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/pb"
	"time"
)

type Status_Checker struct {
	Rbt_Cpid_Time *rbtree.Rbtree
	Rbt_Time_Cpid *rbtree.Rbtree
}

type Cpid_Time_Status struct {
	Cpid     uint64
	RecvTime int64
	CpidTime uint64
}

func (x Cpid_Time_Status) Less(than rbtree.Item) bool {
	return x.Cpid < than.(Cpid_Time_Status).Cpid
}

type Time_Cpid_Status struct {
	RecvTime  int64
	Cpid      uint64
	CpidTime  uint64
	DBID      uint32
	StationID uint32
}

func (x Time_Cpid_Status) Less(than rbtree.Item) bool {
	return x.RecvTime < than.(Time_Cpid_Status).RecvTime
}

var G_Status_Checker *Status_Checker

func GetStatusChecker() *Status_Checker {
	if G_Status_Checker == nil {
		G_Status_Checker = &Status_Checker{
			Rbt_Cpid_Time: rbtree.New(),
			Rbt_Time_Cpid: rbtree.New(),
		}
	}

	return G_Status_Checker
}

func (sc *Status_Checker) Insert(cpid uint64, time_stamp uint64, recv_time_stamp int64, db_id uint32, station_id uint32) {
	//1.insert into rbt_cpid_time
	//    1.1 if has -> update timevalue
	//    1.2 if not has -> insert
	item := sc.Rbt_Cpid_Time.InsertOrGet(Cpid_Time_Status{
		Cpid:     cpid,
		RecvTime: recv_time_stamp,
		CpidTime: time_stamp,
	})
	if item.(Cpid_Time_Status).RecvTime != recv_time_stamp { // means update
		sc.Rbt_Time_Cpid.Delete(Time_Cpid_Status{
			Cpid:     cpid,
			RecvTime: recv_time_stamp,
			CpidTime: time_stamp,
		})
	}
	sc.Rbt_Time_Cpid.Insert(Time_Cpid_Status{
		Cpid:      cpid,
		RecvTime:  recv_time_stamp,
		CpidTime:  time_stamp,
		DBID:      db_id,
		StationID: station_id,
	})

}

func (sc *Status_Checker) Del(recv_time_stamp int64, cpid uint64, time_stamp uint64) {
	//1.del the rbt_time_cpid
	sc.Rbt_Time_Cpid.Delete(Time_Cpid_Status{
		Cpid:     cpid,
		RecvTime: recv_time_stamp,
		CpidTime: time_stamp,
	})

	sc.Rbt_Cpid_Time.Delete(Cpid_Time_Status{
		Cpid:     cpid,
		RecvTime: recv_time_stamp,
		CpidTime: time_stamp,
	})
}

func (sc *Status_Checker) Min() (int64, uint64, uint64, uint32, uint32) {
	time_cpid_status := sc.Rbt_Time_Cpid.Min().(Time_Cpid_Status)

	return time_cpid_status.RecvTime, time_cpid_status.Cpid, time_cpid_status.CpidTime, time_cpid_status.DBID, time_cpid_status.StationID
}

func (sc *Status_Checker) Check() {
	current_time := time.Now().Unix()
	var recv_time int64
	var cpid uint64
	var cp_time_stamp uint64
	var db_id uint32
	var station_id uint32
	for {
		recv_time, cpid, cp_time_stamp, db_id, station_id = sc.Min()
		if current_time-recv_time > int64(conf.GetConf().Redis.OfflineThreshold) {
			GetRedisSocket().ChargingPilesChan <- &Report.ChargingPileStatus{
				Cpid:      cpid,
				Status:    Report.ChargingPileStatus_OFFLINE,
				Id:        db_id,
				StationId: station_id,
				Timestamp: cp_time_stamp,
			}
		} else {
			return
		}
	}
}
