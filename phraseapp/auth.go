package phraseapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bgentry/speakeasy"
)

type AuthHandler struct {
	Username string `cli:"opt --username desc='username used for authentication'"`
	Token    string `cli:"opt --token desc='token used for authentication'"`
	TFA      bool   `cli:"opt --tfa desc='use Two-Factor Authentication'"`
	Host     string "https://api.phraseapp.com/"
	Config   string `cli:"opt --path default=$HOME/.config/phraseapp/config.json desc='path to the config file'"`
}

var authH *AuthHandler

func RegisterAuthHandler(a *AuthHandler) {
	authH = a
}

func (a *AuthHandler) readConfig() error {
	tmpA := new(AuthHandler)

	path := os.ExpandEnv(a.Config)
	_, err := os.Stat(path)
	switch {
	case err == nil: // ignore
		fh, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fh.Close()
		err = json.NewDecoder(fh).Decode(&tmpA)
		if err != nil {
			return err
		}
	case os.IsNotExist(err):
		// ignore
	default:
		return err
	}

	// Set custom 

	// Only set token if username not specified on commandline.
	if tmpA.Token != "" && a.Token == "" && a.Username == "" {
		a.Token = tmpA.Token
	} else if tmpA.Username != "" && a.Username == "" {
		a.Username = tmpA.Username
	}

	if tmpA.TFA && a.Username == "" {
		a.TFA = tmpA.TFA
	}

	// Set custom host if specified
	if tmpA.Host != "" {
		a.Host = tmpA.Host
	}else{
		a.Host = "https://api.phraseapp.com/v2/"
	}

	return nil
}

func authenticate(req *http.Request) error {
	if authH == nil {
		return fmt.Errorf("no auth handler registered")
	}

	if err := authH.readConfig(); err != nil {
		return err
	}

	if err := authH.validate(); err != nil {
		return err
	}

	switch {
	case authH.Token != "":
		req.Header.Set("Authorization", "token "+authH.Token)
	case authH.Username != "":
		pwd, err := speakeasy.Ask("Password: ")
		if err != nil {
			return err
		}
		req.SetBasicAuth(authH.Username, pwd)

		if authH.TFA { // TFA only required for username+password based login.
			token, err := speakeasy.Ask("TFA-Token: ")
			if err != nil {
				return err
			}
			req.Header.Set("X-PhraseApp-OTP", token)
		}
	}

	return nil
}

func (ah *AuthHandler) validate() error {
	switch {
	case ah.Username == "" && ah.Token == "":
		return fmt.Errorf("either username or token must be given")
	default:
		return nil
	}
}
