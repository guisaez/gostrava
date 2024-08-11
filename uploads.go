package gostrava

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// *************** Types ********************

type Upload struct {
	ID         int    `json:"id"`          // The unique identifier of the upload
	IDSrt      string `json:"id_str"`      // The unique identifier of the upload in string format
	ExternalID string `json:"external_id"` // The external identifier of the upload
	ActivityID int    `json:"activity_id"` // The identifier of the activity this upload resulted into
	Error      string `json:"error"`       // The error associated with this upload
	Status     string `json:"status"`      // The status of this upload
}

// *************** Methods ********************

type UploadService service

const uploads string = "/api/v3/uploads"

// CreateUploadRequest represents the parameters for uploading an activity
type CreateUploadRequest struct {
	File        *os.File // File should be of type *os.File
	Name        string
	Description string
	Trainer     string
	Commute     string
	DataType    string
	ExternalID  string
}

// UploadActivity uploads a new data file to create an activity from. Requires activity:write scope.
func (s *UploadService) UploadActivity(ctx context.Context, accessToken string, activity CreateUploadRequest) (*Upload, *http.Response, error) {

	urlStr := uploads

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	upload := new(Upload)
	resp, err := s.client.DoAndParse(ctx, req, upload)
	if err != nil {
		return nil, resp, err
	}

	return upload, resp, err
}

// GetById Returns an upload for a given identifier. Requires activity:write scope.
func (s *UploadService) GetById(ctx context.Context, accessToken string, id int) (*Upload, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d", uploads, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	upload := new(Upload)
	resp, err := s.client.DoAndParse(ctx, req, upload)
	if err != nil {
		return nil, resp, err
	}

	return upload, resp, err
}
