package gostrava

import (
	"context"
	"fmt"
)

type StravaGears struct {
	AccessToken string
	*StravaClient
}

// Returns an equipment using its identifier.
func (sc *StravaGears) GetById(ctx context.Context, id int64) (*DetailedGear, error) {

	path := fmt.Sprintf("/gear/%d", id)

	var resp DetailedGear
	if err := sc.get(ctx, sc.AccessToken, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}