package redis_socket

import (
	"github.com/HuKeping/rbtree"
	"github.com/giskook/charging_pile_tss/conf"
	"github.com/giskook/charging_pile_tss/pb"
	//"log"
	"sync"
	"time"
)

type Status_Checker struct {
	Rbt_Cpid_Time   *rbtree.Rbtree
	Mutex_Cpid_Time sync.Mutex

	Rbt_Time_Cpid   *rbtree.Rbtree
	Mutex_Time_Cpid sync.Mutex
}

type Cpid_Time_Status struct {
	Cpid     uint64
	RecvTime int64
	CpidTime uint64
}

func (x Cpid_Time_Status) Less(than rbtree.Item) bool {
	return x.Cpid < than.(Cpid_Time_Status).Cpid
}

type Cpid_Status struct {
	Cpid      uint64
	CpidTime  uint64
	DBID      uint32
	StationID uint32
}

type Time_Cpid_Status struct {
	RecvTime int64
	Cpids    []*Cpid_Status
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
	sc.Mutex_Cpid_Time.Lock()
	sc.Mutex_Time_Cpid.Lock()
	defer func() {
		sc.Mutex_Cpid_Time.Unlock()
		sc.Mutex_Time_Cpid.Unlock()
	}()
	//1.insert into rbt_cpid_time
	//    1.1 if has ->update timevalue
	//    1.2 if not has -> insert
	sc.Rbt_Cpid_Time.InsertOrGet(Cpid_Time_Status{
		Cpid:     cpid,
		RecvTime: recv_time_stamp,
		CpidTime: time_stamp,
	})
	// 1. if do not have then add
	time_cpid := sc.Rbt_Time_Cpid.Get(Time_Cpid_Status{
		RecvTime: recv_time_stamp,
	})
	if time_cpid == nil {
		sc.Rbt_Time_Cpid.Insert(Time_Cpid_Status{
			RecvTime: recv_time_stamp,
			Cpids: []*Cpid_Status{
				&Cpid_Status{
					Cpid:      cpid,
					CpidTime:  time_stamp,
					DBID:      db_id,
					StationID: station_id,
				},
			},
		})
	} else {
		time_cpid_status := time_cpid.(Time_Cpid_Status)
		time_cpid_status.Cpids = append(time_cpid_status.Cpids, &Cpid_Status{
			Cpid:      cpid,
			CpidTime:  time_stamp,
			DBID:      db_id,
			StationID: station_id,
		})
		sc.Rbt_Time_Cpid.Delete(Time_Cpid_Status{
			RecvTime: recv_time_stamp,
		})

		sc.Rbt_Time_Cpid.Insert(Time_Cpid_Status{
			RecvTime: recv_time_stamp,
			Cpids:    time_cpid_status.Cpids,
		})
	}

}

func (sc *Status_Checker) Del(recv_time_stamp int64) {
	sc.Mutex_Cpid_Time.Lock()
	sc.Mutex_Time_Cpid.Lock()
	defer func() {
		sc.Mutex_Cpid_Time.Unlock()
		sc.Mutex_Time_Cpid.Unlock()
	}()

	//1.del the rbt_time_cpid
	time_cpid_item := sc.Rbt_Time_Cpid.Get(Time_Cpid_Status{
		RecvTime: recv_time_stamp,
	})
	sc.Rbt_Time_Cpid.Delete(Time_Cpid_Status{
		RecvTime: recv_time_stamp,
	})
	cpids_status := time_cpid_item.(Time_Cpid_Status)
	for _, _cpid := range cpids_status.Cpids {
		sc.Rbt_Cpid_Time.Delete(Cpid_Time_Status{
			Cpid: _cpid.Cpid,
		})
	}
}

func (sc *Status_Checker) Min() (int64, []*Cpid_Status) {
	sc.Mutex_Time_Cpid.Lock()
	defer func() {
		sc.Mutex_Time_Cpid.Unlock()
	}()

	if sc.Rbt_Time_Cpid.Len() == 0 {
		return 0, nil
	}

	time_cpid_status := sc.Rbt_Time_Cpid.Min().(Time_Cpid_Status)

	return time_cpid_status.RecvTime, time_cpid_status.Cpids
}

func (sc *Status_Checker) Check() {
	current_time := time.Now().Unix()
	var recv_time int64
	var cpids_status []*Cpid_Status
	for {
		recv_time, cpids_status = sc.Min()
		if cpids_status != nil {
			if current_time-recv_time > int64(conf.GetConf().Redis.OfflineThreshold) {
				for _, cpid := range cpids_status {
					GetRedisSocket().ChargingPilesChan <- &Report.ChargingPileStatus{
						Cpid:      cpid.Cpid,
						Status:    uint32(PROTOCOL_CHARGE_PILE_STATUS_OFFLINE),
						Id:        cpid.DBID,
						StationId: cpid.StationID,
						Timestamp: cpid.CpidTime,
					}
				}

				sc.Del(recv_time)
			} else {
				return
			}
		} else {
			return
		}
	}
}
