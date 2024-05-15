package go_strava

import (
	"encoding/json"
	"net/http"
)

// Parses the ResponseBody into a Fault
func HandleBadResponse(resp *http.Response) error {
	if  resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errorResp Fault
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return &errorResp
		}
	}

	return nil
}