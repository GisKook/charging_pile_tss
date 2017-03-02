package db

import (
	"fmt"
	//"github.com/golang/protobuf/proto"
	"log"
	"time"
)

const (
	TRANS_TABLE_NAME_FMT string = "t_charge_pile_200601"
	SQL_UPDATE_TABLE     string = "UPDATE %s SET start_number=%.2f, end_number=%.2f, electricity=%.2f, money=%.2f, cost=%.2f, time=%d, start_time=time_stamp '%s', end_time=time_stamp '%s', status=%d WHERE order_number = %s"
)

func (db_socket *DbSocket) ProccessTransaction() {
	for {
		select {
		case transactions := <-db_socket.TransactionChan:
			transcation_ids := ""
			tx, err := db_socket.Db.Begin()
			if err != nil {
				log.Println("ProccessTransaction")
				log.Println(err)
			}

			for trans := range transactions {
				update_sql := fmt.Sprintf(SQL_UPDATE_TABLE, GetTableName(trans.StartTime), trans.StartMeterReading, trans.EndMeterReading, trans.ChargingCapacity, trans.ChargingCost, trans.ChargingCostEle, trans.ChargingDuration, GetTime(trans.StartTime), GetTime(trans.EndTime), trans.Status, trans.TransactionID)

				tx.Exec(update_sql)
				transcation_ids += trans.TransactionID + ","
				db_socket.NotifyChan <- transcation_ids
			}
			err = tx.Commit()
			if err != nil {
				log.Println("ProccessTransaction commit error")
				log.Println(err)
			}
		}
	}
}

func GetTableName(t uint64) string {
	_t := time.Unix(int64(t), 0)
	return _t.Format(TRANS_TABLE_NAME_FMT)
}

func GetTime(t uint64) string {
	_t := time.Unix(int64(t), 0)
	return _t.Format("2006-01-02 15:04:05")
}

//func SendNotify(transcation_ids string) {
//	paras := []*Report.Param{
//		&Report.Param{
//			Type:    Report.Param_STRING,
//			Strpara: transcation_ids,
//		},
//	}
//
//	command := &Report.Command{
//		Type:  Report.CMT_NOTIFY_TRANSCATION,
//		Paras: paras,
//	}
//
//	data, _ := proto.Marshal(command)
//
//	GetNsqSocketInstance().Send(conf.GetConf().ProducerConf.TopicWeChat, data)
//}
