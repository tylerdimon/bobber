package http

import (
	"bytes"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordRequestHandler(t *testing.T) {
	mockRequestService := bobber.NewMockRequestService(t)
	mockWebsocketService := bobber.NewMockWebsocketService(t)

	requestToSave := bobber.Request{
		Method:      "POST",
		Host:        "",
		Path:        "/requests/test",
		Body:        `{"some":"json","body":"values"}`,
		NamespaceID: "1234",
		EndpointID:  "4567",
	}

	savedRequest := &bobber.Request{
		ID:          mock.UUIDString,
		Method:      "POST",
		Host:        "",
		Path:        "/requests/test",
		Timestamp:   mock.ParseTime(mock.TimestampString),
		Body:        `{"some":"json","body":"values"}`,
		Headers:     nil,
		NamespaceID: "1234",
		EndpointID:  "4567",
	}

	mockRequestService.EXPECT().Match("POST", "/requests/test").Return("1234", "4567", "a response").Once()
	mockRequestService.EXPECT().Add(requestToSave).Return(savedRequest, nil).Once()
	mockWebsocketService.EXPECT().Broadcast(savedRequest).Once()

	requestBody := bytes.NewBufferString(`{"some":"json","body":"values"}`)
	req, err := http.NewRequest("POST", "/requests/test", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := RequestHandler{
		Service:          mockRequestService,
		WebsocketService: mockWebsocketService,
	}
	handlerFunc := http.HandlerFunc(handler.RecordRequestHandler)
	handlerFunc.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected '%d' but got '%d'", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "a response" {
		t.Errorf("expected '%v' but got '%v'", "Request received", rr.Body.String())
	}
}

// TODO convert to index test
//func TestGetAllRequestsHandler(t *testing.T) {
//	mockRequestService, mockWebsocketService := setup()
//	handler := RequestHandler{
//		Service:          mockRequestService,
//		WebsocketService: mockWebsocketService,
//	}
//
//	expectedRequest1 := bobber.Request{
//		ID:        mock.UUIDString,
//		Method:    "",
//		URL:       "123",
//		Host:      "",
//		Path:      "",
//		Timestamp: mock.TimestampString,
//		Body:      "",
//		Headers:   []bobber.Header{},
//	}
//	expectedRequest2 := bobber.Request{
//		ID:        mock.UUIDString,
//		Method:    "",
//		URL:       "456",
//		Host:      "",
//		Path:      "",
//		Timestamp: mock.TimestampString,
//		Body:      "",
//		Headers:   []bobber.Header{},
//	}
//	_, err := mockRequestService.Add(expectedRequest1)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	_, err = mockRequestService.Add(expectedRequest2)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	req, err := http.NewRequest("GET", "/api/requests/all", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handlerFunc := http.HandlerFunc(handler.RequestIndexHandler)
//
//	handlerFunc.ServeHTTP(rr, req)
//
//	if rr.Code != http.StatusOK {
//		t.Errorf("expected '%d' but got '%d'", http.StatusOK, rr.Code)
//	}
//
//	//expectedRequestStrings := []string{expectedRequest1.String(), expectedRequest2.String()}
//	//expectedBody, err := json.Marshal(expectedRequestStrings)
//	//if err != nil {
//	//	t.Fatal(err)
//	//}
//	//
//	//if rr.Body.String() != string(expectedBody) {
//	//	t.Errorf("expected '%v' but got '%v'", "Request received", rr.Body.String())
//	//}
//}

func DeleteAllRequestsHandler(t *testing.T) {

}
