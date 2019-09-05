package database

type TbSchedule struct {
	ID                 uint
	IntervalForReading int `gorm:"column:intervalforreading"`
	IntervalForEvent   int `gorm:"column:intervalforevent"`
}

func (TbSchedule) TableName() string {
	return "tb_schedule"
}
