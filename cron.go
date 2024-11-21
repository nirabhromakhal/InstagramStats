package main

import (
	_const "InstagramStats/const"
	"InstagramStats/service"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func main() {
	// initialize the service
	db, err := gorm.Open(postgres.Open(_const.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	rapidApiService := service.NewRapidApiService(_const.RapidApiInstagramHost, _const.RapidApiInstagramKey)
	instagramService := service.NewInstagramService(db, rapidApiService)

	// create a scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	// add a job to the scheduler
	_, err = scheduler.NewJob(
		gocron.DurationJob(
			60*time.Second,
		),
		gocron.NewTask(
			instagramService.UpdateAllInstagramChannels,
		),
	)
	if err != nil {
		panic(err)
	}

	// start
	scheduler.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(3 * time.Minute):
	}

	// when you're done, shut it down
	err = scheduler.Shutdown()
	if err != nil {
		panic(err)
	}
}
