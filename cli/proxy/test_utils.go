package proxy

import (
	"net/http"
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

func getClients() ([]*http.Client, error) {
	proxies, err := getTestProxies()
	if err != nil {
		return nil, err
	}
	clients := make([]*http.Client, len(proxies))
	for i, proxy := range proxies {
		client, err := NewProxyClient(proxy)
		if err != nil {
			return nil, err
		}
		clients[i] = client
	}

	return clients, nil
}

func NewTestingClientPool(
	t *testing.T, rateLimit time.Duration,
) *ClientPool {
	clients, err := getClients()
	if err != nil {
		t.Fatal(err)
	}

	pool, err := NewClientPool(rateLimit, clients...)
	if err != nil {
		t.Fatal(err)
	}

	return pool
}
