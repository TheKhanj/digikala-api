package proxy

import (
	"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	proxies, err := getTestProxies()
	if err != nil {
		t.Fatal(err)
	}

	proxy := proxies[0]
	client, err := NewProxyClient(proxy)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf(
			"unexpected response status code (%d)",
			res.StatusCode,
		)
	}
}
