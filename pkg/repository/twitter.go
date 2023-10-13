package repository

import (
	"context"
	"time"

	"github.com/snigle/dddExercise/pkg/domain"
	"github.com/snigle/dddExercise/pkg/repository/connectors"
)

type twitter struct {
	api connectors.HTTPClient
}

func NewTwitterUsername(api connectors.HTTPClient) domain.ITwitterUsername {
	return twitter{api: api}
}

// IsAvailable implements domain.ITwitterUsername.
func (t twitter) IsAvailable(ctx context.Context, username domain.TwitterUsername) (bool, error) {
	availabe := len(username.String()) > 0 && username.String()[0] > 'o'
	time.Sleep(time.Second)
	return availabe, nil
}
