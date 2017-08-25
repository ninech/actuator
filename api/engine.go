package api

import "github.com/gin-gonic/gin"

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
	c.JSON(200, gin.H{"message": "Event received. Thank you."})
}
