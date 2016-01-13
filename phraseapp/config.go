package phraseapp

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type Config struct {
	*Credentials

	ProjectID  string
	Page       *int
	PerPage    *int

	FileFormat string

	Defaults   map[string]map[string]interface{}

	Targets    []byte
	Sources    []byte
}

func (cfg *Config) UnmarshalYAML(unmarshal func(i interface{}) error) error {
	if cfg.Credentials == nil {
		cfg.Credentials = new(Credentials)
	}

	t := struct{ Phraseapp map[string]interface{} }{}
	if err := unmarshal(&t); err != nil {
		return err
	}

	for k, v := range t.Phraseapp {
		switch k {
		// phraseapp.Credentials parameters:
		case "access_token":
			cfg.Credentials.Token = v.(string)
		case "host":
			cfg.Credentials.Host = v.(string)
		case "debug":
			cfg.Credentials.Debug = v.(bool)
		case "username", "tfa":
			return fmt.Errorf("username and tfa not supported in config")
		// ProjectID used if required.
		case "project_id":
			cfg.ProjectID = v.(string)
		case "page":
			page := v.(int)
			cfg.Page = &page
		case "perpage":
			perpage := v.(int)
			cfg.PerPage = &perpage
		case "file_format":
			cfg.FileFormat = v.(string)
		// Special pull and push action configuration.
		case "push":
			var err error
			cfg.Sources, err = yaml.Marshal(v)
			if err != nil {
				return err
			}
		case "pull":
			var err error
			cfg.Targets, err = yaml.Marshal(v)
			if err != nil {
				return err
			}
		// Arbitrary command defaults.
		case "defaults":
			cfg.Defaults = map[string]map[string]interface{}{}
			for path, config := range v.(map[interface{}]interface{}) {
				cfg.Defaults[path.(string)] = map[string]interface{}{}
				for option, value := range config.(map[interface{}]interface{}) {
					cfg.Defaults[path.(string)][option.(string)] = value
				}
			}
		// ignore
		default:
			return fmt.Errorf("unknown key found: %s", k)
		}
	}

	return nil
}
