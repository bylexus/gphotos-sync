package main

import (
	"fmt"
	"io"

	"alexi.ch/gphoto-sync/auth"
)

func main() {
	// The client is prepared with the OAuth token, and will take care of authenticating
	// and refreshing the token by itself.
	client, err := auth.ConfigureHttpClient()
	if err != nil {
		panic(err)
	}

	ret, err := client.Get("https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=20")
	if err != nil {
		panic(err)
	}
	defer ret.Body.Close()
	body, err := io.ReadAll(ret.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
