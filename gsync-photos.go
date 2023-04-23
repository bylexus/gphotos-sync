package main

import (
	"fmt"

	auth "alexi.ch/gphoto-sync/lib"
)

func main() {
	// The client is prepared with the OAuth token, and will take care of authenticating
	// and refreshing the token by itself.
	client, err := auth.ConfigureHttpClient()
	if err != nil {
		panic(err)
	}

	filter := auth.MediaFilter{
		DateFilter: &auth.DateFilter{
			Dates: []auth.Date{{Year: 2022}, {Year: 2023}},
		},
	}

	items, err := auth.LoadMediaItems(client, filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Items read: %d\n", len(*items))
}
