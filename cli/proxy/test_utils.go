package proxy

import (
	"os"
	"path"
	"strings"

	"github.com/thekhanj/digikala-api/cli/internal"
)

func getTestProxies() ([]string, error) {
	bytes, err := os.ReadFile(
		path.Join(
			internal.GetProjectRoot(),
			"test-proxies",
		),
	)
	if err != nil {
		return nil, err
	}

	proxies := strings.Split(string(bytes), "\n")
	ret := make([]string, 0)
	for _, p := range proxies {
		if p != "" {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
