package lib

import (
	"fmt"
	"sync"
)

func NewDownloadManager(config *AppConfig) DownloadManager {
	var wg sync.WaitGroup
	dm := DownloadManager{
		config:      config,
		nrOfThreads: config.NrOfThreads,
		itemQueue:   make(chan MediaItem, config.NrOfThreads),
		waitGroup:   &wg,
	}
	dm.init()
	return dm
}

type DownloadManager struct {
	config      *AppConfig
	nrOfThreads int
	itemQueue   chan MediaItem
	waitGroup   *sync.WaitGroup
}

func (d DownloadManager) init() {
	// Start worker threads:
	for i := 0; i < d.nrOfThreads; i++ {
		d.waitGroup.Add(1)
		fmt.Printf("Starting download worker #%d\n", i+1)
		go d.downloadWorker(i + 1)
	}
}

func (d DownloadManager) Enqueue(item MediaItem) {
	d.itemQueue <- item
}

func (d DownloadManager) DoneEnqueuing() {
	close(d.itemQueue)
}

func (d DownloadManager) downloadWorker(threadNr int) {
	for item := range d.itemQueue {
		destPath, err := item.Download(d.config)
		if err != nil {
			fmt.Printf("   ERROR: %s: %s\n", item.Filename, err)
		} else {
			fmt.Printf("   Thread #%d: download complete: %s\n", threadNr, destPath)
		}
	}
	d.waitGroup.Done()
}

func (d DownloadManager) Wait() {
	d.waitGroup.Wait()
}
