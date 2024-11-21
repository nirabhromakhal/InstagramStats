package controller

import (
	"InstagramStats/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	instagramService *service.InstagramService
}

func NewController(service *service.InstagramService) *Controller {
	return &Controller{instagramService: service}
}

func (ctrl *Controller) LinkInstagram(c *gin.Context) {
	// Get username from request
	username := c.Query("username")
	if username == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "username is required"})
	}

	err := ctrl.instagramService.AddInstagramChannel(username)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Channel added successfully"})
}

func (ctrl *Controller) PollInstagramStats() error {
	err := ctrl.instagramService.UpdateAllInstagramChannels()
	return err
}
