package lib

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

type OauthCodeResponse struct {
	code string
	err  *error
}

func ConfigureHttpClient(config *AppConfig) (*http.Client, error) {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     config.Secrets.GoogleClientId,
		ClientSecret: config.Secrets.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	var err error = nil
	token := config.Secrets.GoogleToken
	if token == nil {
		token, err = createNewToken(conf, ctx)
		if err != nil {
			return nil, err
		}
		config.Secrets.GoogleToken = token
		err = SaveSecrets(config.Secrets)
		if err != nil {
			return nil, err
		}

	} else {
		fmt.Println("Token loaded from token cache")
	}
	client := conf.Client(ctx, token)
	return client, nil
}

func createNewToken(conf *oauth2.Config, ctx context.Context) (*oauth2.Token, error) {
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.

	done := awaitOauthCode(conf)
	response := <-done

	if response.err != nil {
		panic(response.err)
	}

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	code := response.code
	fmt.Println("Got Auth code, proceed with auth")
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func awaitOauthCode(conf *oauth2.Config) chan OauthCodeResponse {
	var srv http.Server
	var listen = "localhost:25123"
	done := make(chan OauthCodeResponse)
	conf.RedirectURL = "http://" + listen + "/"

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	openBrowser(url)
	fmt.Printf("Your browser should open with the Google authentication page.\nPlease follow the instructions there. If the browser does not open by itself, visit the URL for the auth dialog: %v\n\n", url)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got a request")
		params := r.URL.Query()
		code := params["code"]

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><title>gsync photo auth</title></head>
<body>
		<h1>Thank you!</h1>
		<p>You can close this window and go back to your application.</p>
</body>
</html>
		`)
		go func() {
			srv.Shutdown(context.Background())
		}()

		if len(code) > 0 {
			done <- OauthCodeResponse{code: code[0], err: nil}
		} else {
			err := errors.New("auth failed")
			done <- OauthCodeResponse{code: "", err: &err}
		}
		close(done)

	})

	srv.Addr = listen
	go (func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		} else {
			fmt.Println("Server shutdown")
		}
	})()
	return done
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
