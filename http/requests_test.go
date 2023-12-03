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

func setup() (*mock.RequestService, *mock.WebsocketService) {
	mockService := new(mock.RequestService)
	mockService.Gen = mock.Generator()

	websocketService := new(mock.WebsocketService)
	websocketService.Init()

	return mockService, websocketService
}

func TestRecordRequestHandler(t *testing.T) {
	mockRequestService, mockWebsocketService := setup()
	handler := RequestHandler{
		Service:          mockRequestService,
		WebsocketService: mockWebsocketService,
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
	case val := <-mockWebsocketService.Broadcast():
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
		ID:        mock.StaticUUIDValue,
		Method:    "POST",
		URL:       "/requests/test",
		Host:      "",
		Path:      "/requests/test",
		Timestamp: mock.StaticTimeValue,
		Body:      `{"some":"json","body":"values"}`,
		Headers:   "",
	}

	if mockRequestService.Requests[0] != expectedRequest {
		if rr.Body.String() != "Request received" {
			t.Errorf("expected '%v' but got '%v'", expectedRequest, mockRequestService.Requests[0])
		}

	}
}

func TestGetAllRequestsHandler(t *testing.T) {
	mockRequestService, mockWebsocketService := setup()
	handler := RequestHandler{
		Service:          mockRequestService,
		WebsocketService: mockWebsocketService,
	}

	expectedRequest1 := bobber.Request{
		ID:        mock.StaticUUIDValue,
		Method:    "",
		URL:       "123",
		Host:      "",
		Path:      "",
		Timestamp: mock.StaticTimeValue,
		Body:      "",
		Headers:   "",
	}
	expectedRequest2 := bobber.Request{
		ID:        mock.StaticUUIDValue,
		Method:    "",
		URL:       "456",
		Host:      "",
		Path:      "",
		Timestamp: mock.StaticTimeValue,
		Body:      "",
		Headers:   "",
	}
	_, err := mockRequestService.Add(expectedRequest1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = mockRequestService.Add(expectedRequest2)
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
