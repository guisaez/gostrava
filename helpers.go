package gostrava

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleBadResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		if resp.Body != nil {
			defer resp.Body.Close()
			var errResp Error
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
				return err
			}
			return &errResp
		}
		return &Error{Message: "empty response body"}
	}
	return nil
}

func (b *Bounds) toString() string {
	return fmt.Sprintf("%2f,%2f,%2f,%2f", b.SWLat, b.SWLng, b.NELat, b.NELng)
}