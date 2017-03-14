package redis_socket

import (
	"github.com/giskook/charging_pile_tss/db"
	"log"
)

func CalcCost(station_id uint32, pre_cost float32, pre_cost_ele float32, last_meter_reading float32, cur_meter_reading float32, cur_time_stamp uint64) (float32, float32) {
	unit_price, service_price := db.GetDBSocket().GetUnitPrice(station_id, cur_time_stamp)
	log.Println("calc cost")
	log.Printf("station id %d, pre_cost %.2f, pre_cost_ele %.2f, last_meter_reading %.2f, cur_meter_reading %.2f, unit_price %.2f, service_price %.2f", station_id, pre_cost, pre_cost_ele, last_meter_reading, cur_meter_reading, unit_price, service_price)

	log.Println("calc cost")

	cost_total := pre_cost + (cur_meter_reading-last_meter_reading)*(unit_price+service_price)
	cost_ele := pre_cost_ele + (cur_meter_reading-last_meter_reading)*unit_price
	if cost_total < 0 {
		cost_total = 0
	}
	if cost_ele < 0 {
		cost_ele = 0
	}

	return cost_total, cost_ele
}
