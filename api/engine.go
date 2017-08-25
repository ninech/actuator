package api

import "github.com/gin-gonic/gin"

// GetMainEngine constructs the main web engine for the API server
func GetMainEngine() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")

	v1.GET("/health", endpointHealth)

	return r
}

func endpointHealth(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
