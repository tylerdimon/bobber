package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"github.com/tylerdimon/bobber/mocks"
	"github.com/tylerdimon/bobber/static"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordRequestHandler(t *testing.T) {
	mockRequestService := mocks.NewRequestService(t)
	mockWebsocketService := mocks.NewWebsocketService(t)

	requestToSave := bobber.Request{
		Method:      "POST",
		Host:        "",
		Path:        "/requests/test",
		Body:        `{"some":"json","body":"values"}`,
		NamespaceID: "1234",
		EndpointID:  "4567",
	}
	savedRequest := &bobber.Request{}

	mockRequestService.EXPECT().Match("POST", "/requests/test").Return("1234", "4567", "a response").Once()
	mockRequestService.EXPECT().Add(requestToSave).Return(savedRequest, nil).Once()
	mockWebsocketService.EXPECT().Broadcast(savedRequest).Once()

	requestBody := bytes.NewBufferString(`{"some":"json","body":"values"}`)
	req, err := http.NewRequest("POST", "/requests/test", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	handler := RequestHandler{
		Service:          mockRequestService,
		WebsocketService: mockWebsocketService,
	}
	handlerFunc := http.HandlerFunc(handler.RecordRequestHandler)
	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "a response", rr.Body.String())
}

func TestRequestIndexHandler(t *testing.T) {
	static.ParseHTML()

	mockRequestService := mocks.NewRequestService(t)

	request1 := &bobber.Request{
		ID:      mock.UUIDString,
		Method:  "GET",
		Host:    "",
		Path:    "/path",
		Body:    "",
		Headers: []bobber.Header{},
	}
	request2 := &bobber.Request{
		ID:      mock.UUIDString,
		Method:  "POST",
		Host:    "",
		Path:    "/another/path",
		Body:    "",
		Headers: []bobber.Header{},
	}

	mockRequestService.EXPECT().GetAll().Return([]*bobber.Request{request1, request2}, nil).Once()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := RequestHandler{
		Service: mockRequestService,
	}
	handlerFunc := http.HandlerFunc(handler.RequestIndexHandler)
	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)

	// TODO test order, timestamps
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "GET /path")
	assert.Contains(t, rr.Body.String(), "POST /another/path")
}

func TestDeleteAllRequestsHandler(t *testing.T) {
	static.ParseHTML()

	mockRequestService := mocks.NewRequestService(t)

	mockRequestService.EXPECT().DeleteAll().Return(nil).Once()

	req, err := http.NewRequest("DELETE", "/requests", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := RequestHandler{
		Service: mockRequestService,
	}
	handlerFunc := http.HandlerFunc(handler.DeleteAllRequestsHandler)
	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "/", rr.Header().Get("Location"))
}
