package database

type TbGwInfo struct {
	ID         uint
	IPAddr     string
	MacAddr    string
	SocketPort uint
	GatewayID  uint
	FreqID     uint
	CenterFreq uint
	Addr       string `gorm:"column:Addr"`
	FanAmount  uint
}

func (TbGwInfo) TableName() string {
	return "tb_gw_info"
}
