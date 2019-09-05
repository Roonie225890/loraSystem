package database

import "time"

type TbMeterEvent struct {
	ID              uint
	MeterID         string
	CreatedDatetime time.Time `gorm:"column:createdDatetime"`
	EventType       string    `gorm:"column:eventtype"`
	EventCode       uint
	Status          uint8
	DevAddr         uint `gorm:"column:devAddr"`
	ObisIdx         uint
	Attri           uint
}

func (TbMeterEvent) TableName() string {
	return "tb_meter_event"
}
