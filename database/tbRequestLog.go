package database

type TbRequestLog struct {
	ID               uint
	DevAddr          uint `gorm:"column:devAddr"`
	CodingRt         int
	CodingRp         int
	InvokeIDPriority int `gorm:"column:Invoke_Id_Priority"`
	ObisIdx          int
	AttriMethd       int `gorm:"column:Attri_Methd"`
	Flags            int `gorm:"column:Flags"`
	CntOrBlk         int `gorm:"column:Cnt_or_Blk"`
	Choice           int
	ResultOrBlk      int `gorm:"column:Result_or_BlkNum"`
	Conditional      int `gorm:"column:Conditional"`
	RtData           string
	RpData           string
}

func (TbRequestLog) TableName() string {
	return "tb_request_log"
}
