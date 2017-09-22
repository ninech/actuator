package actuator

import (
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
}

// Endpoint defines an interface for all web hook endpoints
type Endpoint interface {
	Handle(c *gin.Context) (int, interface{})
}

// Logger is the central logger for the actuator
var Logger = log.New(os.Stdout, "", log.LstdFlags)

// NewWebhookEngine builds a new webhook engine
// It takes a gin mode as argument
func NewWebhookEngine(mode string) WebhookEngine {
	gin.SetMode(mode)
	return WebhookEngine{}
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
	eventEndpoint := NewEventEndpoint(c.Request)
	code, message := eventEndpoint.Handle()
	c.JSON(code, message)
}
