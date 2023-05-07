package main

import (
	"fmt"

	lib "alexi.ch/gphotos-sync/lib"
	"github.com/jessevdk/go-flags"
)

func main() {
	cmdOpts := lib.CmdOptions{}
	args, err := flags.Parse(&cmdOpts)
	if err != nil {
		panic(err)
	}
	config := lib.CreateAppConfig(args, cmdOpts)

	process(&config)
}

func process(config *lib.AppConfig) {
	// The client is prepared with the OAuth token, and will take care of authenticating
	// and refreshing the token by itself.
	client, err := lib.ConfigureHttpClient(config)
	if err != nil {
		panic(err)
	}

	// We use the mediaItems:search request. This request accepts a filter
	// to limit our items to certain criteria:
	filter := lib.MediaFilter{}

	if len(config.DateFilter) > 0 {
		err := filter.AppendDatesFromStrings(config.DateFilter)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Filters:\n%s\n\n", filter)

	// We load the items:
	// LoadMediaItems returns a channel, which whill get filled with MediaItems,
	// and used as a queue to be processed until closed (all items processed).
	itemsChannel := lib.LoadMediaItems(client, filter)
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
