package proxy

import (
	"testing"
	"time"

	"github.com/thekhanj/digikala-api/cli/config"
)

func getTestProxies() ([]string, error) {
	config, err := config.ReadTestConfig()
	if err != nil {
		return nil, err
	}

	proxies, err := config.Api.Client.GetProxies()
	if err != nil {
		return nil, err
	}

	return proxies, nil
}

func NewTestingClientPool(
	t *testing.T, rateLimit time.Duration,
) *ClientPool {
	proxies, err := getTestProxies()
	if err != nil {
		t.Fatal(err)
	}
	clients, err := NewProxyClientList(proxies)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := NewClientPool(rateLimit, clients...)
	if err != nil {
		t.Fatal(err)
	}

	return pool
}
