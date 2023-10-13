package username

import (
	"context"

	"github.com/snigle/dddExercise/pkg/domain"
)

type Username interface {
	CanUseUsername(ctx context.Context, username string) (bool, error)
}

type username struct {
	twitter domain.ITwitterUsername
}

func NewUsername(twitter domain.ITwitterUsername) Username {
	return username{
		twitter: twitter,
	}
}

func (u username) CanUseUsername(ctx context.Context, username string) (bool, error) {
	twitterUsername, err := domain.NewTwitterUsernameFromString(ctx, username)
	if err != nil {
		return false, err
	}
	return u.twitter.IsAvailable(ctx, twitterUsername)
}
