package database

import "loranet20181205/exception"

type TbEdConf struct {
	ID          uint
	EdID        uint
	Boardsn     string
	DevAddr     int `gorm:"column:devAddr"`
	DRstep      int `gorm:"column:DRstep"`
	RX1DRoffset int `gorm:"column:RX1DRoffset"`
	RX2DataRate int `gorm:"column:RX2DataRate"`
	RX2FC       int `gorm:"column:RX2FC"`
	Delay       int `gorm:"column:Delay"`
	DataRate    int `gorm:"column:DataRate"`
	TXPower     int `gorm:"column:TXPower"`
	RFU         int `gorm:"column:RFU"`
	Linkparam   int
	Uplinkfc    int
}

func (TbEdConf) TableName() string {
	return "tb_ed_conf"
}

func GetEdConf(devAddrI int) TbEdConf {

	var edConf []TbEdConf
	db, err := DbConnect()
	exception.CheckError(err)
	defer db.Close()
	db.Where("devAddr = ?", devAddrI).Find(&edConf)

	return edConf[0]
}
