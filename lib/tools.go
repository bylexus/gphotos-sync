package lib

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kirsle/configdir"
)

func MapToQueryString(input *map[string]string) string {
	var parts []string = make([]string, 0)
	for key, value := range *input {
		parts = append(parts, fmt.Sprintf("%s=%s", url.PathEscape(key), url.PathEscape(value)))
	}
	return strings.Join(parts, "&")
}

func GetAppConfigPath() string {
	configPath := configdir.LocalConfig("gphotos-sync")
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		panic(err)
	}

	return configPath
}
