package database

type TbEdRX2 struct {
	ID   uint
	EdID uint
	DCC  int `gorm:"column:dCC"`
	DFC  int `gorm:"column:dFC"`
	SF   int `gorm:"column:SF"`
}

func (TbEdRX2) TableName() string {
	return "tb_ed_rx2"
}
