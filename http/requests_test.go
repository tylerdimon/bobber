package http

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"github.com/tylerdimon/bobber/mocks"
	"github.com/tylerdimon/bobber/static"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRecordRequestHandler(t *testing.T) {
	mockRequestService := mocks.NewRequestService(t)
	mockWebsocketService := mocks.NewWebsocketService(t)

	s := Server{router: mux.NewRouter()}
	handler := RequestHandler{
		RequestService:   mockRequestService,
		WebsocketService: mockWebsocketService,
	}
	handler.RegisterRequestRoutes(s.router)

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

	rr := httptest.NewRecorder()
	s.serveHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "a response", rr.Body.String())
}

func TestRequestIndexHandler(t *testing.T) {
	static.ParseHTML()

	mockRequestService := mocks.NewRequestService(t)

	s := Server{router: mux.NewRouter()}
	handler := RequestHandler{
		RequestService: mockRequestService,
	}
	handler.RegisterRequestRoutes(s.router)

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

	rr := httptest.NewRecorder()
	s.serveHTTP(rr, req)

	// TODO test order, timestamps
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "GET /path")
	assert.Contains(t, rr.Body.String(), "POST /another/path")
}

func TestDeleteAllRequestsHandler(t *testing.T) {
	static.ParseHTML()

	mockRequestService := mocks.NewRequestService(t)
	s := Server{router: mux.NewRouter()}
	handler := RequestHandler{RequestService: mockRequestService}
	handler.RegisterRequestRoutes(s.router)

	mockRequestService.EXPECT().DeleteAll().Return(nil).Once()

	formData := url.Values{}
	formData.Set("_method", "DELETE")
	req, err := http.NewRequest("POST", "/request", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	s.serveHTTP(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "/", rr.Header().Get("Location"))
}

func TestRequestDetailHandler(t *testing.T) {
	static.ParseHTML()

	rs := mocks.NewRequestService(t)

	s := Server{router: mux.NewRouter()}
	handler := RequestHandler{
		RequestService: rs,
	}
	handler.RegisterRequestRoutes(s.router)

	r := &bobber.Request{
		ID:      mock.UUIDString,
		Method:  "GET",
		Host:    "",
		Path:    "/path",
		Body:    "",
		Headers: []bobber.Header{},
	}

	rs.EXPECT().GetById(mock.UUIDString).Return(r, nil).Once()

	req, err := http.NewRequest("GET", fmt.Sprintf("/request/%s", r.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s.serveHTTP(rr, req)

	// TODO test order, timestamps
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "GET /path")
}
