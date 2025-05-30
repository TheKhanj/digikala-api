package config

import (
	"os"
)

func ReadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = c.UnmarshalJSON(bytes)
	if err != nil {
		return nil, err
	}

	return c, nil
}
