package gostrava

import (
	"context"
	"fmt"
)

type StravaUploads struct {
    AccessToken string
	*StravaClient
}

// Returns an upload for a given identifier. Requires activity:write scope.
func (sc *StravaUploads) GetUpload(ctx context.Context, id int64) (*Upload, error) {
    
    path := fmt.Sprintf("/uploads/%d", id)
    var resp Upload
    if err := sc.get(ctx, sc.AccessToken, path, nil, &resp); err != nil {
        return nil, err
    }

    return &resp, nil
}