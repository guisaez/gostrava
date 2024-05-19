package gostrava

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type StravaStreams baseModule

type StreamKey string

var AllowedStreamKeys = struct {
	Time           StreamKey
	Distance       StreamKey
	LatLng         StreamKey
	Altitude       StreamKey
	VelocitySmooth StreamKey
	Heartrate      StreamKey
	Cadence        StreamKey
	Watts          StreamKey
	Temp           StreamKey
	Moving         StreamKey
	GradeSmooth    StreamKey
}{
	"time",
	"distance",
	"latlng",
	"altitude",
	"velocity_smooth",
	"heartrate",
	"cadence",
	"watts",
	"temp",
	"moving",
	"grade_smooth",
}

// Returns the given activity's streams. Requires activity:read scope. Requires activity:read_all scope for Only Me activities.
func (sc *StravaStreams) GetActivityStreams(ctx context.Context, access_token string, id int64, keys []StreamKey, key_by_type bool) (*StreamSet, error) {
    path := fmt.Sprintf("/activities/%d/streams", id)

 
    queryParams := url.Values{}
    queryParams.Set("key_by_type", strconv.FormatBool(key_by_type))

    if len(keys) > 0 {
        var keysString []string
        for _, k := range keys {
            keysString = append(keysString, string(k))
        }
        queryParams.Set("keys", strings.Join(keysString, ","))
    }

    var resp StreamSet
    err := sc.client.get(ctx, access_token, path, queryParams, &resp)
    if err != nil {
        return nil, err
    }

    return &resp, nil
}

// Returns the given route's streams. Requires read_all scope for private routes.
func (sc *StravaStreams) GetRouteStreams(ctx context.Context, access_token string, id int64) (*StreamSet, error) {
	path := fmt.Sprintf("/routes/%d/streams", id)

	var resp StreamSet
    err := sc.client.get(ctx, access_token, path, nil, &resp)
    if err != nil {
        return nil, err
    }

    return &resp, nil
}

// Returns a set of streams for a segment effort completed by the authenticated athlete. Requires read_all scope.
func (sc *StravaStreams) GetSegmentEffortStreams(ctx context.Context, access_token string, id int64, keys []StreamKey, key_by_type bool) (*StreamSet, error) {
    path := fmt.Sprintf("/segment_efforts/%d/streams", id)

    queryParams := url.Values{}
    queryParams.Set("key_by_type", strconv.FormatBool(key_by_type))

    if len(keys) > 0 {
        var keysString []string
        for _, k := range keys {
            keysString = append(keysString, string(k))
        }
        queryParams.Set("keys", strings.Join(keysString, ","))
    }

    var resp StreamSet
    err := sc.client.get(ctx, access_token, path, queryParams, &resp)
    if err != nil {
        return nil, err
    }

    return &resp, nil
}

// Returns the given segment's streams. Requires read_all scope for private segments.
func (sc *StravaStreams) GetSegmentStreams(ctx context.Context, access_token string, id int64, keys []StreamKey, key_by_type bool) (*StreamSet, error) {
    
	path := fmt.Sprintf("/segments/%d/streams", id)

    queryParams := url.Values{}
    queryParams.Set("key_by_type", strconv.FormatBool(key_by_type))

    if len(keys) > 0 {
        var keysString []string
        for _, k := range keys {
            keysString = append(keysString, string(k))
        }
        queryParams.Set("keys", strings.Join(keysString, ","))
    }

    var resp StreamSet
    err := sc.client.get(ctx, access_token, path, queryParams, &resp)
    if err != nil {
        return nil, err
    }

    return &resp, nil
}