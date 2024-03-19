package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func gormSetup() *gorm.DB {
	dsn := "root:172983456@tcp(localhost:3306)/goapi?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Album{})
	//	db.Create(&Album{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99})
	//	db.Create(&Album{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99})
	//	db.Create(&Album{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99})
	return db
}
