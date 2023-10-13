package main

import (
	"context"
	"fmt"
	"log"

	"github.com/snigle/dddExercise/pkg/repository"
	"github.com/snigle/dddExercise/pkg/repository/connectors"
	"github.com/snigle/dddExercise/pkg/usecase/username"
)

func main() {
	ctx := context.Background()
	usernameInput := "todo"

	twitterConnector, err := connectors.NewHTTPClient(ctx, "https://twitter.com/")
	if err != nil {
		log.Fatalf("fail: %s", err)
	}

	valid, err := username.NewUsername(repository.NewTwitterUsername(twitterConnector)).CanUseUsername(ctx, usernameInput)
	if err != nil {
		log.Fatalf("fail: %s", err)
	}

	if valid {
		fmt.Printf("username %s is valid\n", usernameInput)
	} else {
		fmt.Printf("username %s is not valid\n", usernameInput)
	}
}
