package db

import (
	"database/sql"
	"github.com/giskook/charging_pile_tss/base"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

func (db_socket *DbSocket) LoadAllPrices() error {
	st, err := db_socket.Db.Prepare("SELECT id, charge_station_id, start_time, end_time, electricity_price, service_price FROM t_electricity_price")

	if err != nil {
		return err
	}

	r, e := st.Query()
	defer st.Close()

	if e != nil {
		return e
	}

	var sql_id sql.NullInt64
	var sql_station_id sql.NullInt64
	var sql_start_time pq.NullTime
	var sql_end_time pq.NullTime
	var sql_electricity_price sql.NullFloat64
	var sql_service_price sql.NullFloat64

	db_socket.ChargingPricesMutex.Lock()

	for r.Next() {
		err = r.Scan(
			&sql_id,
			&sql_station_id,
			&sql_start_time,
			&sql_end_time,
			&sql_electricity_price,
			&sql_service_price,
		)
		id := uint64(GetInt64Value(sql_id, 0))
		station_id := uint64(GetInt64Value(sql_station_id, 0))
		start_time, _ := GetTimeValue(sql_start_time)
		end_time, _ := GetTimeValue(sql_end_time)
		electricity_price := float32(GetFloat64Value(sql_electricity_price, 0.0))
		service_price := float32(GetFloat64Value(sql_service_price, 0))

		if err != nil {
			log.Println(err.Error())
			db_socket.ChargingPricesMutex.Unlock()
			return err
		}
		db_socket.ChargingPrices[station_id] =
			append(db_socket.ChargingPrices[station_id], &base.ChargingPrice{
				ID:              id,
				Start_hour:      uint8(start_time.Hour()),
				Start_min:       uint8(start_time.Minute()),
				End_hour:        uint8(end_time.Hour()),
				End_min:         uint8(end_time.Minute()),
				Elec_unit_price: electricity_price,
				Service_price:   service_price,
			})
	}
	db_socket.ChargingPricesMutex.Unlock()

	defer r.Close()

	return nil
}

func (db_socket *DbSocket) parse_payload_price(notify string) {
	switch notify[0] {
	case 'U':
		db_socket.update_price(notify)
	case 'I':
		db_socket.insert_price(notify)
	case 'D':
		db_socket.del_price(notify)
	}

}

func (db_socket *DbSocket) parse_payload_price_common(payload string) (uint64, uint64, *base.ChargingPrice) {
	values := strings.Split(payload, "^")
	id, _ := strconv.ParseUint(values[1], 10, 64)
	station_id, _ := strconv.ParseUint(values[2], 10, 64)
	start_time_string := values[4]
	end_time_string := values[5]
	elec_unit_price, _ := strconv.ParseFloat(values[6], 32)
	service_price, _ := strconv.ParseFloat(values[7], 32)
	start_time, _ := time.Parse("15:04:05", start_time_string)
	end_time, _ := time.Parse("15:04:05", end_time_string)

	return id, station_id, &base.ChargingPrice{
		ID:              id,
		Start_hour:      uint8(start_time.Hour()),
		Start_min:       uint8(start_time.Minute()),
		End_hour:        uint8(end_time.Hour()),
		End_min:         uint8(end_time.Minute()),
		Elec_unit_price: float32(elec_unit_price),
		Service_price:   float32(service_price),
	}
}

func (db_socket *DbSocket) insert_price(payload string) {
	db_socket.ChargingPricesMutex.Lock()
	defer db_socket.ChargingPricesMutex.Unlock()

	_, station_id, charging_price := db_socket.parse_payload_price_common(payload)
	log.Println(db_socket.ChargingPrices[station_id])
	db_socket.ChargingPrices[station_id] = append(db_socket.ChargingPrices[station_id], charging_price)
	log.Println(db_socket.ChargingPrices[station_id])
}

func (db_socket *DbSocket) del_price(payload string) {
	db_socket.ChargingPricesMutex.Lock()
	defer db_socket.ChargingPricesMutex.Unlock()

	id, station_id, _ := db_socket.parse_payload_price_common(payload)
	for i, p := range db_socket.ChargingPrices[station_id] {
		if p.ID == id && len(db_socket.ChargingPrices[station_id]) > 0 {
			db_socket.ChargingPrices[station_id][i] = db_socket.ChargingPrices[station_id][len(db_socket.ChargingPrices[station_id])-1]
			db_socket.ChargingPrices[station_id][len(db_socket.ChargingPrices[station_id])-1] = nil
			db_socket.ChargingPrices[station_id] = db_socket.ChargingPrices[station_id][:len(db_socket.ChargingPrices[station_id])-1]
			return
		}
	}
}

func (db_socket *DbSocket) update_price(payload string) {
	db_socket.ChargingPricesMutex.Lock()
	defer db_socket.ChargingPricesMutex.Unlock()

	id, station_id, price := db_socket.parse_payload_price_common(payload)
	for i, p := range db_socket.ChargingPrices[station_id] {
		if p.ID == id {
			db_socket.ChargingPrices[station_id][i] = price
			return
		}
	}
}

func (db_socket *DbSocket) GetUnitPrice(station_id uint32, cur_time uint64) (float32, float32) {
	db_socket.ChargingPricesMutex.Lock()
	defer db_socket.ChargingPricesMutex.Unlock()

	charging_prices := db_socket.ChargingPrices[uint64(station_id)]
	if charging_prices != nil {
		//_rest_time := (cur_time % (24 * 60 * 60))
		//_hour := uint8(_rest_time / (60 * 60))
		//_min := uint8((_rest_time % (60 * 60)) / 60)
		__cur_time := time.Unix(int64(cur_time), 0)
		_cur_time := __cur_time.Local()
		_hour := uint8(_cur_time.Hour())
		_min := uint8(_cur_time.Minute())
		log.Println("llll")
		log.Println(_hour)
		log.Println(_min)

		for index := range charging_prices {
			price := charging_prices[index]
			if base.In_period(price.Start_hour, price.Start_min,
				price.End_hour, price.End_min, _hour, _min) {

				return price.Elec_unit_price, price.Service_price
			}
		}
	}
	log.Println("there is no price here GetUnitPrice")

	return 0.0, 0.0
}
