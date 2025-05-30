package config

import (
	"path"

	"github.com/thekhanj/digikala-api/cli/internal"
)

func ReadTestConfig() (*Config, error) {
	return ReadConfig(
		path.Join(internal.GetProjectRoot(), "github-config.json"),
	)
}
