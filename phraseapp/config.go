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

	var ok bool
	cfgErrStr := "configuration key %q has invalid value\nsee https://phraseapp.com/docs/developers/cli/configuration/"
	for k, v := range t.Phraseapp {
		switch k {
		// phraseapp.Credentials parameters:
		case "access_token":
			if cfg.Credentials.Token, ok = v.(string); !ok {
				return fmt.Errorf(cfgErrStr, k)
			}
		case "host":
			if cfg.Credentials.Host, ok = v.(string); !ok {
				return fmt.Errorf(cfgErrStr, k)
			}
		case "debug":
			if cfg.Credentials.Debug, ok = v.(bool); !ok {
				return fmt.Errorf(cfgErrStr, k)
			}
		case "username", "tfa":
			return fmt.Errorf("username and tfa not supported in config")
		// ProjectID used if required.
		case "project_id":
			if cfg.ProjectID, ok = v.(string); !ok {
				return fmt.Errorf(cfgErrStr, k)
			}
		case "page":
			page, ok := v.(int)
			if !ok {
				return fmt.Errorf(cfgErrStr, k)
			}
			cfg.Page = &page
		case "perpage":
			perpage, ok := v.(int)
			if !ok {
				return fmt.Errorf(cfgErrStr, k)
			}
			cfg.PerPage = &perpage
		case "file_format":
			if cfg.FileFormat, ok = v.(string); !ok {
				return fmt.Errorf(cfgErrStr, k)
			}
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
			val, ok := v.(map[interface{}]interface{})
			if !ok { return fmt.Errorf(cfgErrStr, k) }

			cfg.Defaults = map[string]map[string]interface{}{}
			for rawPath, rawConfig := range val {
				path, ok := rawPath.(string)
				if !ok { return fmt.Errorf(cfgErrStr, rawPath) }

				config, ok := rawConfig.(map[interface{}]interface{})
				if !ok { return fmt.Errorf(cfgErrStr, rawPath)}

				cfg.Defaults[path] = map[string]interface{}{}
				for rawKey, value := range config {
					key, ok := rawKey.(string)
					if !ok { return fmt.Errorf(cfgErrStr, key) }

					cfg.Defaults[path][key] = value
				}
			}
		// ignore
		default:
			return fmt.Errorf("configuration key %q invalid\nsee https://phraseapp.com/docs/developers/cli/configuration/", k)
		}
	}

	return nil
}
