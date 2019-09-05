package database

import "time"

type TbEdStatus struct {
	ID         uint
	EdID       uint
	RpDatetime time.Time `gorm:"column:rp_datetime"`
}

func (TbEdStatus) TableName() string {
	return "tb_ed_status"
}
