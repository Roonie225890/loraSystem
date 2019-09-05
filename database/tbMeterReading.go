package database

type TbMeterReading struct {

	// gorm.Model
	ID           uint
	MeterID      string
	Recordnum    int
	ReadingType1 string
	Value1       float64
	ReadingType2 string
	Value2       float64
	Ts           string
	Status       int8
	ReadingType  string
	Value        float64
}

func (TbMeterReading) TableName() string {
	return "tb_meter_reading"
}

// func main() {
// 	db, err := gorm.Open("mysql", "root:harbinger@/iot_db?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	var reading []tbMeterReading
// 	// db.AutoMigrate(&Product{})
//
// 	// db.Create(&tb_meter_reading{Meter_id: "22222"})
// 	// db.Model(&product).Where("price = ?", 0).Update("price", 5000)
// 	// db.First(&reading, 1) // 查询id为1的product
//
// 	// db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
// 	db.Where("recordnum = ?", 4).Find(&reading)
// 	fmt.Println(reading)
// 	defer db.Close()
// }
