package main

import (
	"context"
	"fmt"
	"github.com/k-yomo/moodle"
	"net/url"
)

func main() {
	ctx := context.Background()
	serviceURL, err := url.Parse("https://my.uopeople.edu")
	if err != nil {
		panic(err)
	}
	moodleClient, err := moodle.NewClientWithLogin(ctx, serviceURL, &moodle.LoginParams{
		Username: "S000000",
		Password: "password",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", moodleClient)
}
