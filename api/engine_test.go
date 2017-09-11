package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ninech/actuator/api"
)

func ServeRequest(method string, endpoint string) *httptest.ResponseRecorder {
	server := api.GetMainEngine()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(method, endpoint, nil)

	server.ServeHTTP(recorder, request)

	return recorder
}

func TestUnknownRoute(t *testing.T) {
	recorder := ServeRequest("GET", "/yolo")

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, recorder.Code, recorder.HeaderMap.Get("Location"))
	}
}

func TestEndpointHealth(t *testing.T) {
	recorder := ServeRequest("GET", "/v1/health")

	if recorder.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusOK, recorder.Code, recorder.HeaderMap.Get("Location"))
	}
}
