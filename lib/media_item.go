package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type MediaItem struct {
	Id            string
	Description   string
	ProductUrl    string
	BaseUrl       string
	MimeType      string
	Filename      string
	MediaMetadata MediaMetadata
}

// / Downloads a media item to a certain path
func (m *MediaItem) Download(config *AppConfig) error {
	var downloadPath = m.mediaDownloadPath(config)
	var downloadUrl string

	if m.MediaMetadata.Photo != nil {
		downloadUrl = m.BaseUrl + "=d"
	} else if m.MediaMetadata.Video != nil {
		if m.MediaMetadata.Video.Status == "READY" {
			downloadUrl = m.BaseUrl + "=vd"
		} else {
			return errors.New("video not ready")
		}
	}

	// create destination path
	err := os.MkdirAll(downloadPath, 0755)
	if err != nil {
		return err
	}

	// check destination file:
	// if it exists, and its creation date is the same as the
	// file to download, skip it:
	fullpath := filepath.Join(downloadPath, m.Filename)
	if !m.shouldOverride(fullpath, config) {
		return nil
	}

	// create the download request:
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read download, write file:
	fmt.Printf("Downloading %s to %s\n", m.Filename, fullpath)
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	// change the file's times to the file's creation time
	mktime := m.MediaMetadata.GetCreationTime()
	os.Chtimes(fullpath, mktime, mktime)

	return nil
}

func (m *MediaItem) mediaDownloadPath(config *AppConfig) string {
	var path = config.BaseOutputPath
	t, err := time.Parse(time.RFC3339, m.MediaMetadata.CreationTime)
	if err != nil {
		path = filepath.Join(path, "UNKNOWN")
	} else {
		path = filepath.Join(path, fmt.Sprintf("%d", t.Year()))
	}
	return path
}

type MediaMetadata struct {
	CreationTime string
	Height       string
	Width        string
	Photo        *PhotoMetadata
	Video        *VideoMetadata
}

func (m *MediaMetadata) GetCreationTime() time.Time {
	t, err := time.Parse(time.RFC3339, m.CreationTime)
	if err != nil {
		return time.Now()
	} else {
		return t
	}
}

type PhotoMetadata struct {
	CameraMake      string
	CameraModel     string
	FocalLength     float64
	ApertureFNumber float64
	IsoEquivalent   int
	ExposureTime    string
}

type VideoMetadata struct {
	CameraMake  string
	CameraModel string
	Fps         float64
	Status      string
}

type SearchRequestBody struct {
	PageSize  int         `json:"pageSize"`
	PageToken *string     `json:"pageToken,omitempty"`
	Filters   MediaFilter `json:"filters"`
}

type MediaItemResponse struct {
	Item    *MediaItem
	BatchNr int
	Err     error
}
type MediaItemsChannel chan MediaItemResponse

func LoadMediaItems(client *http.Client, filter MediaFilter) MediaItemsChannel {
	var pageSize = 100
	var channel = make(MediaItemsChannel, pageSize)
	var pageToken *string = nil
	var batchNr = 0

	// process items in a separate thread, and use the channel as buffer
	go func() {
		defer close(channel)

		// process in batches of pageSize items
		for {
			batchNr += 1

			var requestBody = SearchRequestBody{
				PageSize:  pageSize,
				PageToken: nil,
				Filters:   filter,
			}

			// the pageToken is a "next page" pointer for the Google API.
			// If the last batch call returned a nextPageToken, we attach
			// it to the next request:
			if pageToken != nil {
				requestBody.PageToken = pageToken
			}

			// Form the request
			requestBodyString, err := json.Marshal(requestBody)
			if err != nil {
				channel <- MediaItemResponse{Item: nil, BatchNr: batchNr, Err: err}
				return
			}

			// execute the request:
			ret, err := client.Post(
				"https://photoslibrary.googleapis.com/v1/mediaItems:search",
				"application/json",
				bytes.NewReader(requestBodyString),
			)

			if err != nil {
				channel <- MediaItemResponse{Item: nil, BatchNr: batchNr, Err: err}
				return
			}

			// process the response:
			body, err := io.ReadAll(ret.Body)
			ret.Body.Close()
			if err != nil {
				channel <- MediaItemResponse{Item: nil, BatchNr: batchNr, Err: err}
				return
			}

			mediaItems := MediaItemsResponse{}

			err = json.Unmarshal(body, &mediaItems)
			if err != nil {
				channel <- MediaItemResponse{Item: nil, BatchNr: batchNr, Err: err}
				return
			}

			if len(mediaItems.MediaItems) > 0 {
				// yield returned MediaItems over the channel:
				for _, item := range mediaItems.MediaItems {
					// make a copy of item here: item is re-used in every
					// loop, making &item point to the SAME variable on each loop.
					// We need to copy it here before returing a pointer:
					resItem := item
					channel <- MediaItemResponse{
						Item:    &resItem,
						BatchNr: batchNr,
						Err:     nil,
					}
				}
			}
			// are we done? Yes if the response did not sent a nextPageToken:
			if mediaItems.NextPageToken == nil {
				break
			} else {
				pageToken = mediaItems.NextPageToken
			}
		}
	}()

	return channel
}

func (m *MediaItem) shouldOverride(destFile string, config *AppConfig) bool {
	if config.ForceOverride {
		return true
	}

	// check destination file:
	// if it exists, and its mod date >= the remote file, skip the download
	info, err := os.Stat(destFile)
	if err == nil {
		if info.Mode().IsRegular() && !config.ForceNewerOverride {
			fmt.Printf("Skipping %s: local file exists\n", m.Filename)
			return false
		}
		if info.ModTime().Compare(m.MediaMetadata.GetCreationTime()) == 0 {
			fmt.Printf("Skipping %s: local file has same timestamp\n", m.Filename)
			return false
		} else if info.ModTime().Compare(m.MediaMetadata.GetCreationTime()) > 0 {
			fmt.Printf("Skipping %s: local file is newer\n", m.Filename)
			return false
		}
	}
	return true
}
