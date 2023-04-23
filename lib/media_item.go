package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MediaItem struct {
	Id         string
	ProductUrl string
	MimeType   string
	Filename   string
}

type SearchRequestBody struct {
	PageSize  int         `json:"pageSize"`
	PageToken *string     `json:"pageToken"`
	Filters   MediaFilter `json:"filters"`
}

func LoadMediaItems(client *http.Client, filter MediaFilter) (*[]MediaItem, error) {
	var result = make([]MediaItem, 0)
	var pageToken *string = nil
	var requestBody = SearchRequestBody{
		PageSize:  100,
		PageToken: nil,
		Filters:   filter,
	}
	var pageSize = 100
	var itemCounter = 0

	for {
		fmt.Printf("Working on items %d - %d...\n", itemCounter+1, itemCounter+pageSize)
		if pageToken != nil {
			requestBody.PageToken = pageToken
		} else {
			requestBody.PageToken = nil
		}
		itemCounter = itemCounter + pageSize

		requestBodyString, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		ret, err := client.Post(
			"https://photoslibrary.googleapis.com/v1/mediaItems:search",
			"application/json",
			bytes.NewReader(requestBodyString),
		)

		if err != nil {
			return nil, err
		}
		defer ret.Body.Close()

		body, err := io.ReadAll(ret.Body)
		if err != nil {
			return nil, err
		}

		mediaItems := MediaItemsResponse{}

		err = json.Unmarshal(body, &mediaItems)
		if err != nil {
			return nil, err
		}

		if len(mediaItems.MediaItems) > 0 {
			result = append(result, mediaItems.MediaItems...)
		}
		if mediaItems.NextPageToken == nil {
			pageToken = nil
			break
		} else {
			pageToken = mediaItems.NextPageToken
		}
	}

	return &result, nil
}
