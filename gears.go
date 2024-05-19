package gostrava

import (
	"context"
	"fmt"
)

type StravaGears baseModule

// Returns an equipment using its identifier.
func (sc *StravaGears) GetById(ctx context.Context, access_token string, id int64, ) (*DetailedGear, error) {

	path := fmt.Sprintf("/gear/%d", id)

	var resp DetailedGear
	if err := sc.client.get(ctx, access_token, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}