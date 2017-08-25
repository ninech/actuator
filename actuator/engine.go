package actuator

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

// GetMainEngine constructs the main web engine for the API server
func GetMainEngine() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/health", endpointHealth)
		v1.POST("/event", endpointEvent)
	}
	return r
}

func endpointHealth(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func endpointEvent(c *gin.Context) {
	var event PullRequestEvent

	if c.BindJSON(&event) == nil {
		message := fmt.Sprintf("Event for pull request #%d received. Thank you.", event.Number)
		logger.Println("Received event.")
		c.JSON(200, gin.H{"message": message})
	} else {
		c.JSON(400, gin.H{"message": "Invalid JSON payload."})
	}
}
