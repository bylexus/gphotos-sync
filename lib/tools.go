package lib

import (
	"fmt"
	"net/url"
	"strings"
)

func MapToQueryString(input *map[string]string) string {
	var parts []string = make([]string, 0)
	for key, value := range *input {
		parts = append(parts, fmt.Sprintf("%s=%s", url.PathEscape(key), url.PathEscape(value)))
	}
	return strings.Join(parts, "&")
}
