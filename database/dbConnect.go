package database

import "github.com/jinzhu/gorm"

func DbConnect() (*gorm.DB, error) {

	db, err := gorm.Open("mysql", "root:pwd@/iot_db?charset=utf8&parseTime=True&loc=Local")

	return db, err
}
