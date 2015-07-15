package phraseapp

import (
	"fmt"
	"net/http"

	"github.com/bgentry/speakeasy"
)

type DefaultParams Action
type Action map[string]map[string]interface{}

type AuthCredentials struct {
	Username string `cli:"opt --username desc='username used for authentication'"`
	Token    string `cli:"opt --token desc='token used for authentication'"`
	TFA      bool   `cli:"opt --tfa desc='use Two-Factor Authentication'"`
	Host     string `cli:"opt --host desc='Host to send Request to'"`
	Debug    bool   `cli:"opt --verbose desc='Verbose output'"`
}

var authC *AuthCredentials

func RegisterAuthCredentials(cmdAuth *AuthCredentials, defaultCredentials *AuthCredentials) {
	authC = new(AuthCredentials)

	if cmdAuth.Token != "" && authC.Token == "" && authC.Username == "" {
		authC.Token = cmdAuth.Token
	} else if cmdAuth.Username != "" && authC.Username == "" {
		authC.Username = cmdAuth.Username
	}

	if cmdAuth.TFA && authC.Username == "" {
		authC.TFA = cmdAuth.TFA
	}

	if cmdAuth.Debug == true || ((defaultCredentials != nil) && defaultCredentials.Debug == true) {
		EnableDebug()
	}

	if cmdAuth.Host != "" {
		authC.Host = cmdAuth.Host
	} else {
		if defaultCredentials != nil && defaultCredentials.Host != "" {
			authC.Host = defaultCredentials.Host
		}
	}

	if authC.Host == "" {
		authC.Host = "https://api.phraseapp.com"
	}

	notSet := authC.Token == "" && authC.Username == ""
	if notSet && defaultCredentials.Token != "" && defaultCredentials != nil {
		authC.Token = defaultCredentials.Token
	}
	if notSet && defaultCredentials != nil && defaultCredentials.Username != "" {
		authC.Username = defaultCredentials.Username
	}
}

func authenticate(req *http.Request) error {
	if authC == nil {
		return fmt.Errorf("no auth handler registered")
	}

	if err := authC.validate(); err != nil {
		return err
	}

	switch {
	case authC.Token != "":
		req.Header.Set("Authorization", "token "+authC.Token)
	case authC.Username != "":
		pwd, err := speakeasy.Ask("Password: ")
		if err != nil {
			return err
		}
		req.SetBasicAuth(authC.Username, pwd)

		if authC.TFA { // TFA only required for username+password based login.
			token, err := speakeasy.Ask("TFA-Token: ")
			if err != nil {
				return err
			}
			req.Header.Set("X-PhraseApp-OTP", token)
		}
	}

	return nil
}

func (ah *AuthCredentials) validate() error {
	switch {
	case ah.Username == "" && ah.Token == "":
		return fmt.Errorf("either username or token must be given")
	default:
		return nil
	}
}
