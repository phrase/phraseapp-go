package phraseapp

import (
	"fmt"
	"net/http"

	"os"

	"github.com/bgentry/speakeasy"
)

type DefaultParams Action
type Action map[string]map[string]interface{}

type Client struct {
	Credentials *Credentials
}

type Credentials struct {
	Username string `cli:"opt --username desc='username used for authentication'"`
	Token    string `cli:"opt --access-token desc='access token used for authentication'"`
	TFA      bool   `cli:"opt --tfa desc='use Two-Factor Authentication'"`
	Host     string `cli:"opt --host desc='Host to send Request to'"`
	Debug    bool   `cli:"opt --verbose desc='Verbose output'"`
}

func NewClient(credentials Credentials, defaultCredentials *Credentials) (*Client, error) {
	client := &Client{Credentials: &Credentials{}}

	envToken := os.Getenv("PHRASEAPP_ACCESS_TOKEN")

	if credentials.Token != "" && client.Credentials.Token == "" && client.Credentials.Username == "" {
		client.Credentials.Token = credentials.Token
	} else if credentials.Username != "" && client.Credentials.Username == "" {
		client.Credentials.Username = credentials.Username
	} else if envToken != "" && credentials.Token == "" && credentials.Username == "" && client.Credentials.Username == "" {
		client.Credentials.Token = envToken
	}

	if credentials.TFA && client.Credentials.Username == "" {
		client.Credentials.TFA = credentials.TFA
	}

	if credentials.Debug == true || ((defaultCredentials != nil) && defaultCredentials.Debug == true) {
		EnableDebug()
	}

	if credentials.Host != "" {
		client.Credentials.Host = credentials.Host
	} else {
		if defaultCredentials != nil && defaultCredentials.Host != "" {
			client.Credentials.Host = defaultCredentials.Host
		}
	}

	if client.Credentials.Host == "" {
		client.Credentials.Host = "https://api.phraseapp.com"
	}

	notSet := client.Credentials.Token == "" && client.Credentials.Username == ""
	if notSet && defaultCredentials != nil && defaultCredentials.Token != "" {
		client.Credentials.Token = defaultCredentials.Token
	}
	if notSet && defaultCredentials != nil && defaultCredentials.Username != "" {
		client.Credentials.Username = defaultCredentials.Username
	}

	return client, nil
}

func (client *Client) authenticate(req *http.Request) error {
	if client.Credentials == nil {
		return fmt.Errorf("no auth handler registered")
	}

	if err := client.Credentials.validate(); err != nil {
		return err
	}

	req.Header.Set("User-Agent", GetUserAgent())
	switch {
	case client.Credentials.Token != "":
		req.Header.Set("Authorization", "token "+client.Credentials.Token)
	case client.Credentials.Username != "":
		pwd, err := speakeasy.Ask("Password: ")
		if err != nil {
			return err
		}
		req.SetBasicAuth(client.Credentials.Username, pwd)

		if client.Credentials.TFA { // TFA only required for username+password based login.
			token, err := speakeasy.Ask("TFA-Token: ")
			if err != nil {
				return err
			}
			req.Header.Set("X-PhraseApp-OTP", token)
		}
	}

	return nil
}

func (ah *Credentials) validate() error {
	switch {
	case ah.Username == "" && ah.Token == "":
		return fmt.Errorf("either username or token must be given")
	default:
		return nil
	}
}
