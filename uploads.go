package go_strava

import (
	"context"
	"fmt"
)

// Returns an upload for a given identifier. Requires activity:write scope.
func (sc *StravaClient) GetUpload(ctx context.Context, id int64) (*Upload, error) {
    
    path := fmt.Sprintf("/uploads/%d", id)
    var resp Upload
    if err := sc.get(ctx, path, nil, &resp); err != nil {
        return nil, err
    }

    return &resp, nil
}
