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

	var err error
	for k, v := range t.Phraseapp {
		switch k {
		// phraseapp.Credentials parameters:
		case "access_token":
			if cfg.Credentials.Token, err = ValidateIsString(k, v); err != nil {
				return err
			}
		case "host":
			if cfg.Credentials.Host, err = ValidateIsString(k, v); err != nil {
				return err
			}
		case "debug":
			if cfg.Credentials.Debug, err = ValidateIsBool(k, v); err != nil {
				return err
			}
		case "username", "tfa":
			return fmt.Errorf("username and tfa not supported in config")
		// ProjectID used if required.
		case "project_id":
			if cfg.ProjectID, err = ValidateIsString(k, v); err != nil {
				return err
			}
		case "page":
			page, err := ValidateIsInt(k, v)
			if err != nil {
				return err
			}
			cfg.Page = &page
		case "perpage":
			perpage, err := ValidateIsInt(k, v)
			if err != nil {
				return err
			}
			cfg.PerPage = &perpage
		case "file_format":
			if cfg.FileFormat, err = ValidateIsString(k, v); err != nil {
				return err
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
			val, err := ValidateIsRawMap(k, v)
			if err != nil {
				return err
			}

			cfg.Defaults = map[string]map[string]interface{}{}
			for path, rawConfig := range val {
				cfg.Defaults[path], err = ValidateIsRawMap("defaults." + path, rawConfig)
				if err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("configuration key %q invalid\nsee https://phraseapp.com/docs/developers/cli/configuration/", k)
		}
	}

	return nil
}

const cfgValueErrStr = "configuration key %q has invalid value\nsee https://phraseapp.com/docs/developers/cli/configuration/"
const cfgKeyErrStr = "configuration key %q has invalid type\nsee https://phraseapp.com/docs/developers/cli/configuration/"
const cfgInvalidKeyErrStr = "configuration key %q unknown\nsee https://phraseapp.com/docs/developers/cli/configuration/"

func ValidateIsString(k string, v interface{}) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf(cfgValueErrStr, k)
	}
	return s, nil
}

func ValidateIsBool(k string, v interface{}) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf(cfgValueErrStr, k)
	}
	return b, nil
}

func ValidateIsInt(k string, v interface{}) (int, error) {
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf(cfgValueErrStr, k)
	}
	return i, nil
}

func ValidateIsRawMap(k string, v interface{}) (map[string]interface{}, error) {
	raw, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf(cfgValueErrStr, k)
	}

	ps := map[string]interface{}{}
	for mk, mv := range raw {
		s, ok := mk.(string)
		if !ok {
			return nil, fmt.Errorf(cfgKeyErrStr, fmt.Sprintf("%s.%v", k, mk))
		}
		ps[s] = mv
	}
	return ps, nil
}

func ParseYAMLToMap(unmarshal func(interface{}) error, keysWithType map[string]interface{})  error {
	m := map[string]interface{}{}
	if err := unmarshal(m); err != nil {
		return err
	}

	var err error
	for k, v := range m {
		value, found := keysWithType[k]
		if !found {
			return fmt.Errorf(cfgInvalidKeyErrStr, k)
		}

		switch val := value.(type) {
		case *string:
			*val, err = ValidateIsString(k, v)
		case *int:
			*val, err = ValidateIsInt(k, v)
		case *bool:
			*val, err = ValidateIsBool(k, v)
		case map[string]interface{}:
			val, err = ValidateIsRawMap(k, v)
		default:
			err = fmt.Errorf(cfgValueErrStr, k)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
