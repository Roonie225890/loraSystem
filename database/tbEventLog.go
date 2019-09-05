package database

type TbEventLog struct {
	ID        uint
	DtStr     string
	ObisIndex uint8 `gorm:"column:OBIS_index"`
	Attribute uint8 `gorm:"column:Attribute"`
	EventCode uint8
}

func (TbEventLog) TableName() string {
	return "tb_event_log"
}
