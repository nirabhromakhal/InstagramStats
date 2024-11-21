package main

import (
	_const "InstagramStats/const"
	"InstagramStats/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(_const.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&entity.User{}, &entity.Channel{}, &entity.ChannelMember{}, &entity.Video{})
	if err != nil {
		panic("failed to migrate database")
	}
}
