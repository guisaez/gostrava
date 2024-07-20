package gostrava

import (
	"net/http"
	"os"
	"strconv"
)

type UploadService service

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

type Upload struct {
	ID         int    `json:"id"`          // The unique identifier of the upload
	IdSrt      string `json:"id_str"`      // The unique identifier of the upload in string format
	ExternalID string `json:"external_id"` // The external identifier of the upload
	ActivityID int    `json:"activity_id"` // The identifier of the activity this upload resulted into
	Error      string `json:"error"`       // The error associated with this upload
	Status     string `json:"status"`      // The status of this upload
}

// Uploads a new data file to create an activity from. Requires activity:write scope.
func (s *UploadService) UploadActivity(accessToken string, data CreateUploadRequest) (*Upload, error) {
	req, err := s.client.NewRequest(RequestOpts{
		Path:        "uploads",
		Method:      http.MethodPost,
		AccessToken: accessToken,
		Body:        data,
	})
	if err != nil {
		return nil, err
	}

	resp := new(Upload)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns an upload for a given identifier. Requires activity:write scope.
func (s *UploadService) GetById(accessToken string, uploadID int) (*Upload, error) {
	req, err := s.client.NewRequest(RequestOpts{
		Path:        "uploads/" + strconv.Itoa(uploadID),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(Upload)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
