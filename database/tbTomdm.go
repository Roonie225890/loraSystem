package database

type TbTomdm struct {
	ID    uint
	MsgID string
	Code  string
}

func (TbTomdm) TableName() string {
	return "tb_tomdm"
}
