package http

import (
	"bytes"
	"github.com/tylerdimon/bobber/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRequestHandler(t *testing.T) {
	mockService := new(mock.RequestService)
	websocketService := new(mock.WebsocketService)
	websocketService.Init()
	handler := RequestHandler{
		Service:          mockService,
		WebsocketService: websocketService,
	}

	// consume channel to prevent blocking
	go func() {
		for range websocketService.Broadcast() {
		}
	}()

	requestBody := bytes.NewBufferString(`{"method":"GET","url":"http://example.com"}`)

	req, err := http.NewRequest("POST", "/add", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.AddRequestHandler)

	handlerFunc.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected '%d' but got '%d'", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "Request received" {
		t.Errorf("expected '%v' but got '%v'", "Request received", rr.Body.String())
	}
}
