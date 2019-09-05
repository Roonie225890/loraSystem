package database

import (
	"log"
	"loranet20181205/exception"
)

type TbMeterList struct {
	NO            uint   `gorm:"column:no"`
	MeterUniqueID string `gorm:"column:meterUniqueID"`
	MeterID       string `gorm:"column:meterID"`
	CustomerID    string `gorm:"column:customerId"`
	AK            string `gorm:"column:AK"`
	GUK           string `gorm:"column:GUK"`
}

func (TbMeterList) TableName() string {
	return "tb_meter_list"
}

func GetMeterList(MeterID string) TbMeterList {

	var meterList []TbMeterList
	db, err := DbConnect()
	exception.CheckError(err)
	defer db.Close()
	db.Where("meterID = ?", MeterID).Find(&meterList)
	if len(meterList) == 0 {
		log.Print("ERROR: Invalid MeterID")
	}
	// log.Print("AK:", meterList[0].AK, " GUK:", meterList[0].GUK)
	return meterList[0]
}
