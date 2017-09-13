package actuator_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ninech/actuator/actuator"
	"github.com/stretchr/testify/assert"
)

func ServeRequest(method string, endpoint string, body string) *httptest.ResponseRecorder {
	engine := actuator.NewWebhookEngine(gin.TestMode)
	router := engine.GetRouter()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(method, endpoint, strings.NewReader(body))

	router.ServeHTTP(response, request)

	return response
}

func TestUnknownRoute(t *testing.T) {
	recorder := ServeRequest("GET", "/yolo", "")

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, recorder.Code, recorder.HeaderMap.Get("Location"))
	}
}

func TestEndpointHealth(t *testing.T) {
	recorder := ServeRequest("GET", "/v1/health", "")

	assert.Equal(t, http.StatusOK, recorder.Code, "they should be equal")
}

func TestEndpointEvent(t *testing.T) {
	response := ServeRequest("POST", "/v1/event", `{"number":1,"action":"opened"}`)
	body, _ := ioutil.ReadAll(response.Body)

	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	assert.JSONEq(t, `{"message":"Event for pull request #1 received. Thank you."}`, string(body), "should be equal")
}

func TestEndpointEventInvalidJSON(t *testing.T) {
	response := ServeRequest("POST", "/v1/event", `{}`)
	body, _ := ioutil.ReadAll(response.Body)

	assert.Equal(t, http.StatusBadRequest, response.Code, "they should be equal")
	assert.JSONEq(t, `{"message":"Invalid JSON payload."}`, string(body), "should be equal")
}
