package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

type Secrets struct {
	GoogleClientId     string `json:"googleClientId"`
	GoogleClientSecret string `json:"googleClientSecret"`
	GoogleToken        *oauth2.Token
}

func LoadSecrets() *Secrets {
	path := GetAppConfigPath()
	configFile := filepath.Join(path, "secrets.json")
	var secrets Secrets

	// file does not exist: create new one
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		secrets = Secrets{}
	} else {
		// file exists, parse:
		fh, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		secrets = Secrets{}
		decoder := json.NewDecoder(fh)
		decoder.Decode(&secrets)
	}

	return &secrets
}

func SaveSecrets(secrets *Secrets) error {
	path := GetAppConfigPath()
	configFile := filepath.Join(path, "secrets.json")

	fh, err := os.Create(configFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	return encoder.Encode(secrets)
}

func (s *Secrets) EnsureUserSecrets() error {
	if len(s.GoogleClientId) == 0 || len(s.GoogleClientSecret) == 0 {
		var clientId string
		var clientSecret string

		fmt.Println("Enter your Google API credentials:")
		fmt.Print("Google Client ID: ")
		_, err := fmt.Scanln(&clientId)
		if err != nil {
			return err
		}

		fmt.Print("Google Client Secret: ")
		_, err = fmt.Scanln(&clientSecret)
		if err != nil {
			return err
		}
		s.GoogleClientId = clientId
		s.GoogleClientSecret = clientSecret
		return SaveSecrets(s)

	}
	return nil
}
