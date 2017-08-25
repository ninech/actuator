package frauschultz_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/frauschultz"
)

func ServeRequest(method string, endpoint string) *httptest.ResponseRecorder {
	server := frauschultz.GetMainEngine()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(method, endpoint, nil)

	server.ServeHTTP(response, request)

	return response
}

func TestUnknownRoute(t *testing.T) {
	recorder := ServeRequest("GET", "/yolo")

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, recorder.Code, recorder.HeaderMap.Get("Location"))
	}
}

func TestEndpointHealth(t *testing.T) {
	recorder := ServeRequest("GET", "/v1/health")

	assert.Equal(t, recorder.Code, http.StatusOK, "they should be equal")
}

func TestEndpointEvent(t *testing.T) {
	response := ServeRequest("POST", "/v1/event")

	assert.Equal(t, response.Code, http.StatusOK, "they should be equal")
}

func TestEndpointEventBody(t *testing.T) {
	response := ServeRequest("POST", "/v1/event")

	body, _ := ioutil.ReadAll(response.Body)

	assert.JSONEq(t, `{"message":"Event received. Thank you."}`, string(body), "should be equal")
}
