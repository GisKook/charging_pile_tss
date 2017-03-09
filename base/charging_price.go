package base

import ()

type TransactionDetail struct {
	TransactionID string
	Status        uint8

	ChargingDuration uint32
	ChargingCapacity float32
	ChargingCost     float32
	ChargingCostEle  float32

	StartTime         uint64
	EndTime           uint64
	StartMeterReading float32
	EndMeterReading   float32
}

type ChargingPrice struct {
	ID              uint64
	Start_hour      uint8
	Start_min       uint8
	End_hour        uint8
	End_min         uint8
	Elec_unit_price float32
	Service_price   float32
}

type ChargingCost struct {
	Uuid string
	Tid  uint64
	Cost uint32
}

type StopChargingNotify struct {
	Uuid string
	Tid  uint64
}

func In_period(start_hour uint8, start_minute uint8, end_hour uint8, end_minute uint8, target_hour uint8, target_minute uint8) bool {
	start_minutes := uint16(start_hour)*60 + uint16(start_minute)
	end_minutes := uint16(end_hour)*60 + uint16(end_minute)
	target_minutes := uint16(target_hour)*60 + uint16(target_minute)

	end_minutes_ := end_minutes - start_minutes
	if end_minutes_ < 0 {
		end_minutes_ += 24 * 60
	}
	target_minutes_ := target_minutes - start_minutes
	if target_minutes_ < 0 {
		target_minutes_ += 24 * 60
	}

	if target_minutes_ > end_minutes_ {
		return false
	}

	return true
}
