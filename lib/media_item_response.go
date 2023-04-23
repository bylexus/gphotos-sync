package lib

type MediaItemsResponse struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken *string     `json:"nextPageToken,omitempty"`
}
