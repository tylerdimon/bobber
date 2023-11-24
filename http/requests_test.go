package http

import (
	"bytes"
	"encoding/json"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRecordRequestHandler(t *testing.T) {
	mockService := new(mock.RequestService)
	websocketService := new(mock.WebsocketService)
	websocketService.Init()
	handler := RequestHandler{
		Service:          mockService,
		WebsocketService: websocketService,
	}

	requestBody := bytes.NewBufferString(`{"some":"json","body":"values"}`)

	req, err := http.NewRequest("POST", "/requests/test", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.RecordRequestHandler)

	go handlerFunc.ServeHTTP(rr, req)

	select {
	case val := <-websocketService.Broadcast():
		// TODO validate message value
		t.Logf("Got message: %s", val)
	case <-time.After(time.Second * 1):
		t.Error("Expected a value to be sent to the channel, but timed out")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("expected '%d' but got '%d'", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "Request received" {
		t.Errorf("expected '%v' but got '%v'", "Request received", rr.Body.String())
	}

	expectedRequest := bobber.Request{
		ID:        "",
		Method:    "POST",
		URL:       "/requests/test",
		Host:      "",
		Path:      "/requests/test",
		Timestamp: "",
		Body:      `{"some":"json","body":"values"}`,
		Headers:   "",
	}

	assertRequestsEqual(t, expectedRequest, mockService.Requests[0])
}

func TestGetAllRequestsHandler(t *testing.T) {
	mockService := new(mock.RequestService)
	websocketService := new(mock.WebsocketService)
	websocketService.Init()
	handler := RequestHandler{
		Service:          mockService,
		WebsocketService: websocketService,
	}

	expectedRequest1 := bobber.Request{
		ID:        "",
		Method:    "",
		URL:       "123",
		Host:      "",
		Path:      "",
		Timestamp: time.Now().String(),
		Body:      "",
		Headers:   "",
	}
	expectedRequest2 := bobber.Request{
		ID:        "",
		Method:    "",
		URL:       "456",
		Host:      "",
		Path:      "",
		Timestamp: time.Now().String(),
		Body:      "",
		Headers:   "",
	}
	_, err := mockService.Add(expectedRequest1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = mockService.Add(expectedRequest2)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/api/requests/all", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.GetAllRequests)

	handlerFunc.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected '%d' but got '%d'", http.StatusOK, rr.Code)
	}

	expectedRequestStrings := []string{expectedRequest1.String(), expectedRequest2.String()}
	expectedBody, err := json.Marshal(expectedRequestStrings)
	if err != nil {
		t.Fatal(err)
	}

	if rr.Body.String() != string(expectedBody) {
		t.Errorf("expected '%v' but got '%v'", "Request received", rr.Body.String())
	}
}

func DeleteAllRequestsHandler(t *testing.T) {

}

func assertRequestsEqual(t *testing.T, expected, actual bobber.Request) {
	// TODO need to mock out IDs and Timestamps then can get rid of this
	// and just compare structs
	if actual.Method != expected.Method ||
		actual.URL != expected.URL ||
		actual.Host != expected.Host ||
		actual.Path != expected.Path ||
		actual.Body != expected.Body ||
		actual.Headers != expected.Headers {
		t.Errorf("expected '%v' but got '%v'", expected, actual)
	}

}
