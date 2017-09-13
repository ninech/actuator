package actuator_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
)

type ServeRequestHeader map[string]string
type ServeRequestOptions struct {
	Method   string
	Endpoint string
	Body     string
	Header   ServeRequestHeader
}

func ServeRequest(options ServeRequestOptions) *httptest.ResponseRecorder {
	engine := actuator.NewWebhookEngine(gin.TestMode)
	router := engine.GetRouter()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(options.Method, options.Endpoint, strings.NewReader(options.Body))

	for key, value := range options.Header {
		request.Header.Add(key, value)
	}

	router.ServeHTTP(response, request)

	return response
}

func TestUnknownRoute(t *testing.T) {
	options := ServeRequestOptions{Method: "GET", Endpoint: "/yolo", Body: ""}
	recorder := ServeRequest(options)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, recorder.Code, recorder.HeaderMap.Get("Location"))
	}
}

func TestEndpointHealth(t *testing.T) {
	options := ServeRequestOptions{Method: "GET", Endpoint: "/v1/health"}
	recorder := ServeRequest(options)

	assert.Equal(t, http.StatusOK, recorder.Code, "they should be equal")
}
