package main

import (
	"fmt"
	"os"

	lib "alexi.ch/gphotos-sync/lib"
	"github.com/jessevdk/go-flags"
)

func main() {
	cmdOpts := lib.CmdOptions{}
	args, err := flags.Parse(&cmdOpts)
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		} else {
			panic(err)
		}
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
	if len(config.DateRangeFilter) > 0 {
		err := filter.AppendDateRangesFromStrings(config.DateRangeFilter)
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
	dm := lib.NewDownloadManager(config)

	for itemResponse := range itemsChannel {
		// copy itemResponse, as it is re-used on each loop run:
		workingItem := itemResponse
		if workingItem.Err != nil {
			fmt.Printf("ERROR: %v\n", workingItem.Err)
		} else {
			// process item
			counter += 1
			// fmt.Printf("Enqueuing item #%d: %s\n", counter, workingItem.Item.Filename)
			dm.Enqueue(*workingItem.Item)
		}
	}
	dm.DoneEnqueuing()
	dm.Wait()
	fmt.Printf("Items processed: %d\n", counter)
}
