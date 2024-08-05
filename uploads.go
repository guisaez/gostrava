package gostrava

import (
	"os"
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
	IDSrt      string `json:"id_str"`      // The unique identifier of the upload in string format
	ExternalID string `json:"external_id"` // The external identifier of the upload
	ActivityID int    `json:"activity_id"` // The identifier of the activity this upload resulted into
	Error      string `json:"error"`       // The error associated with this upload
	Status     string `json:"status"`      // The status of this upload
}
