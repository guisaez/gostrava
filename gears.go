package go_strava

import (
	"context"
	"fmt"
)

// Returns an equipment using its identifier.
func (sc *StravaClient) GetEquipment(ctx context.Context, id int64) (*DetailedGear, error) {

	path := fmt.Sprintf("/gear/%d", id)

	var resp DetailedGear
	if err := sc.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}