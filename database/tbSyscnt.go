package database

import (
	"log"
	"loranet20181205/exception"
	"strconv"
)

type TbSyscnt struct {
	ID          uint
	DevAddr     int `gorm:"column:devAddr"`
	DownlinkCnt int `gorm:"column:downlink_cnt"`
	UplinkCnt   int `gorm:"column:uplink_cnt"`
	InvokeID    int `gorm:"column:Invoke_Id"`
}

func (TbSyscnt) TableName() string {
	return "tb_syscnt"
}

func GetSyscn(devAddrI int) TbSyscnt {
	// var NwKSKey string
	// var AppSKey string
	var downlinkCnt int
	var uplinkCnt int
	var invokeID int
	// var uplink_cnt int
	var syscnt []TbSyscnt
	// Create Connection
	db, err := DbConnect()
	exception.CheckError(err)
	// Close Connection
	defer db.Close()
	// Query
	// qstr3 := "SELECT NwKSKey, AppSKey FROM tb_ed_info WHERE devAddr = " + strconv.Itoa(devAddrI)

	// log.Print(len(edInfo[0].NwkSKey))

	db.Where("devAddr = ?", strconv.Itoa(devAddrI)).Find(&syscnt)

	if len(syscnt) == 0 {
		downlinkCnt = 1
		invokeID = 0
		db.Create(&TbSyscnt{DevAddr: devAddrI, DownlinkCnt: downlinkCnt, UplinkCnt: uplinkCnt, InvokeID: invokeID})
		db.Last(&syscnt)
		log.Printf("Dnlnk: %d Ivk:%d\n", downlinkCnt, invokeID)

		// log.Print("ERROR: Syscnt")
		log.Printf("insertid:%d\n", syscnt)
	} else {
		downlinkCnt = syscnt[0].DownlinkCnt
		invokeID = syscnt[0].InvokeID
		downlinkCnt++
		invokeID++
		rowsAffect := db.Model(&syscnt).Where("devAddr = ?", devAddrI).Updates(map[string]interface{}{"downlink_cnt": downlinkCnt, "Invoke_Id": invokeID}).RowsAffected
		log.Printf("Dnlnk: %d Ivk:%d\n", downlinkCnt, invokeID)
		log.Printf("update:%d\n", rowsAffect)

	}

	return syscnt[0]
}
