package domain

import "context"

type ITwitterUsername interface {
	IsAvailable(ctx context.Context, username TwitterUsername) (bool, error)
}

type TwitterUsername struct {
	value string
}

func (t TwitterUsername) String() string {
	return t.value
}

func NewTwitterUsernameFromString(ctx context.Context, input string) (TwitterUsername, error) {
	resp := TwitterUsername{value: input}
	return resp, nil
}
