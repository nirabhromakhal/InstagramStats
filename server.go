package main

import (
	_const "InstagramStats/const"
	"InstagramStats/controller"
	"InstagramStats/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "OK"})
}

func InitController() *controller.Controller {
	db, err := gorm.Open(postgres.Open(_const.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	rapidApiService := service.NewRapidApiService(_const.RapidApiInstagramHost, _const.RapidApiInstagramKey)
	instagramService := service.NewInstagramService(db, rapidApiService)
	return controller.NewController(instagramService)
}

func main() {
	controller := InitController()
	router := gin.Default()

	// routes
	router.GET("/link/instagram", controller.LinkInstagram)
	router.GET("/health", health)
	router.Run("localhost:8080")
}
