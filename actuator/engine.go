package actuator

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	// DebugMode sets gin into debug mode
	DebugMode string = gin.DebugMode
	// ReleaseMode disables gin debug features
	ReleaseMode string = gin.ReleaseMode
	// TestMode runs gin in test mode
	TestMode string = gin.TestMode
)

// WebhookEngine represents the main struct for the webhook application
type WebhookEngine struct {
	Logger *log.Logger
}

// NewWebhookEngine builds a new webhook engine
// It takes a gin mode as argument
func NewWebhookEngine(mode string) WebhookEngine {
	gin.SetMode(mode)
	return WebhookEngine{
		Logger: log.New(os.Stdout, "", log.LstdFlags)}
}

// GetRouter constructs the main web engine for the API server
func (e *WebhookEngine) GetRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/health", e.endpointHealth)
		v1.POST("/event", e.endpointEvent)
	}
	return r
}

func (e *WebhookEngine) endpointHealth(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func (e *WebhookEngine) endpointEvent(c *gin.Context) {
	var event PullRequestEvent

	if c.BindJSON(&event) == nil {
		message := fmt.Sprintf("Event for pull request #%d received. Thank you.", event.Number)
		e.Logger.Println("Received event.")
		c.JSON(200, gin.H{"message": message})
	} else {
		c.JSON(400, gin.H{"message": "Invalid JSON payload."})
	}
}
