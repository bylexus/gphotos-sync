package main

import (
	"fmt"
	"os"

	auth "alexi.ch/gphoto-sync/lib"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := auth.AppConfig{BaseOutputPath: cwd}
	process(&config)
}

func process(config *auth.AppConfig) {
	// The client is prepared with the OAuth token, and will take care of authenticating
	// and refreshing the token by itself.
	client, err := auth.ConfigureHttpClient()
	if err != nil {
		panic(err)
	}

	// We use the mediaItems:search request. This request accepts a filter
	// to limit our items to certain criteria:
	filter := auth.MediaFilter{
		DateFilter: &auth.DateFilter{
			Dates: []auth.Date{
				// {Year: 2022},
				{Year: 2023}},
		},
	}

	// We load the items:
	// LoadMediaItems returns a channel, which whill get filled with MediaItems,
	// and used as a queue to be processed until closed (all items processed).
	itemsChannel := auth.LoadMediaItems(client, filter)
	counter := 0

	for itemResponse := range itemsChannel {
		counter += 1
		fmt.Printf("Working on Batch %d, item %d...\n", itemResponse.BatchNr, counter)
		if itemResponse.Err != nil {
			fmt.Printf("ERROR: %v\n", itemResponse.Err)
		} else {
			// process item
			fmt.Printf("   Filename: %s\n", itemResponse.Item.Filename)
			itemResponse.Item.Download(config)
		}
	}
	fmt.Printf("Items processed: %d\n", counter)
}
