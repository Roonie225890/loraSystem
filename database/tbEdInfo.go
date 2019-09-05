package database

import (
	"log"
	"loranet20181205/exception"
)

type TbEdInfo struct {
	ID            int
	GwID          uint
	HopBDA        uint   `gorm:"column:hop_BDA"`
	DevAddr       uint   `gorm:"column:devAddr"`
	GwMacAddr     string `gorm:"column:gwMacAddr"`
	NO            uint   `gorm:"column:NO"`
	MeterUniqueID string `gorm:"column:meterUniqueID"`
	MeterID       string `gorm:"column:MeterID"`
	CustomerID    string `gorm:"column:CustomerId"`
	DevEUI        string `gorm:"column:DevEUI"`
	AppEUI        string `gorm:"column:AppEUI"`
	AppKey        string `gorm:"column:AppKey"`
	NwkSKey       string `gorm:"column:NwkSKey"`
	AppSKey       string `gorm:"column:AppSKey"`
	AppNonce      int    `gorm:"column:AppNonce"`
	GUK           string `gorm:"column:GUK"`
	DevNonce      int    `gorm:"column:DevNonce"`
	AK            string `gorm:"column:AK"`
	NetID         int    `gorm:"column:NetID"`
	ActiveStatus  uint   `gorm:"column:activeStatus"`
	Boardsn       string `gorm:"column:board_sn"`
}

func (TbEdInfo) TableName() string {
	return "tb_ed_info"
}

func GetEdInfo(DevEUI string) TbEdInfo {

	var edInfo []TbEdInfo
	db, err := DbConnect()
	exception.CheckError(err)
	defer db.Close()
	db.Where("DevEUI = ?", DevEUI).Find(&edInfo)
	log.Print("DevEUI:", edInfo[0].DevEUI, " AppKey:", edInfo[0].AppKey)
	return edInfo[0]
}
