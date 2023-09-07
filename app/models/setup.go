package models

import (
	"fmt"

	"api_sotr/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DbSotr *gorm.DB

func ConnectDatabase() {
	conf := config.New()
	dsn := conf.Employee.UserName + ":" + conf.Employee.Password + "@tcp(" + conf.Employee.Host + ")/" + conf.Employee.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("db err: (Init) ", err)
	}
	DbSotr = db
}
