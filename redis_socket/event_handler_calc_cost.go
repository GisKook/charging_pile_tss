package redis_socket

import (
	"github.com/giskook/charging_pile_tss/db"
)

func CalcCost(station_id uint32, pre_cost float32, pre_cost_ele float32, last_meter_reading float32, cur_meter_reading float32, cur_time_stamp uint64) (float32, float32) {
	unit_price, service_price := db.GetDBSocket().GetUnitPrice(station_id, cur_time_stamp)

	cost_total := pre_cost + (cur_meter_reading-last_meter_reading)*(unit_price+service_price)
	cost_ele := pre_cost_ele + (cur_meter_reading-last_meter_reading)*unit_price

	return cost_total, cost_ele
}
