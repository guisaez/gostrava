package gostrava

import (
	"context"
	"fmt"
)

type StravaUploads baseModule

// Returns an upload for a given identifier. Requires activity:write scope.
func (sc *StravaUploads) GetUpload(ctx context.Context, access_token string, id int64) (*Upload, error) {
    
    path := fmt.Sprintf("/uploads/%d", id)
    var resp Upload
    if err := sc.client.get(ctx, access_token, path, nil, &resp); err != nil {
        return nil, err
    }

    return &resp, nil
}