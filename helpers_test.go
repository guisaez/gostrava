package gostrava

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// Returns nil if the response 200 < statusCode <= 400
func TestHandleBadResponse(t *testing.T){

	// OK Response
	resp1 := &http.Response{
		StatusCode: http.StatusOK,
		Status: "200 OK",
	}

	err := handleBadResponse(resp1)
	if err != nil {
		t.Errorf("expected no error")
	}

	// Bad Request but Empty Body
	resp2 := &http.Response{
		StatusCode: http.StatusBadRequest,
		Status: "400 BadRequest",
	}

	err = handleBadResponse(resp2)
	if err == nil {
		t.Errorf("expected and error")
	}
}


func TestHandleBadResponseWithBody(t *testing.T){

	errorPayload := &Error{
		Message: "Authorization Error",
		Errors: []ErrorContent{
			{
				Code: "invalid",
				Field: "access_token",
				Resource: "Athlete",
			},
		},
	}

	jsonM, err := json.Marshal(errorPayload)
	if err != nil {
		t.Errorf("error marshalling payload - not related to test case")
	}
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Status: "400 Bad Request",
		Body: io.NopCloser(bytes.NewBuffer(jsonM)),
	}

	err = handleBadResponse(resp)
	if err == nil {
		t.Errorf("expected an error")
	}

	if err.Error() != string(jsonM) {
		t.Errorf("expected error message %s, got %s", string(jsonM), err.Error())
	}
}