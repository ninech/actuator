package frauschultz

import (
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
	logger.Println("Received event.")
	c.JSON(200, gin.H{"message": "Event received. Thank you."})
}
