package proxy

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func getIp(pool *ClientPool) (string, error) {
	req, err := http.NewRequest("GET", "https://ifconfig.io", nil)
	if err != nil {
		return "", err
	}

	res, err := pool.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	s := string(bytes)
	return strings.TrimSpace(s), nil
}

func TestClientPool(t *testing.T) {
	pool := NewTestingClientPool(t, 0)
	defer func() {
		stopped := pool.Shutdown()
		<-stopped
	}()

	ips := make(map[string]bool)
	for i := range pool.clients {
		ip, err := getIp(pool)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%d'th proxy recieved ip %s", i, ip)
		ips[ip] = true
	}

	if len(ips) != len(pool.clients) {
		t.Fatalf(
			"expected %d different ips for each proxy got %d",
			len(pool.clients), len(ips),
		)
	}
}

func TestClientPoolRateLimit(t *testing.T) {
	// By keeping the time between requests high enough, we are making the time
	// needed for each request to complete negligable, so we won't get into
	// "too long" test errors, hopefully.
	rateLimit := time.Second * 3
	pool := NewTestingClientPool(t, rateLimit)
	defer func() {
		stopped := pool.Shutdown()
		<-stopped
	}()

	var wg sync.WaitGroup
	before := time.Now().UnixMilli()

	eachProxyReqCount := 5
	for i := 0; i < eachProxyReqCount; i++ {
		for range pool.clients {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ip, err := getIp(pool)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf(
					"recieved ip %s at %dms",
					ip, time.Now().UnixMilli()-before,
				)
			}()
		}
	}
	wg.Wait()
	after := time.Now().UnixMilli()

	totalTime := int(after - before)
	tooShort := totalTime < int(rateLimit.Milliseconds())*(eachProxyReqCount-1)
	if tooShort {
		t.Fatalf("requests took too short to complete (%d seconds)", totalTime)
	}
	tooLong := totalTime > int(rateLimit.Milliseconds())*(eachProxyReqCount-1)*len(clients)
	if tooLong {
		t.Fatalf("requests took too long to complete (%d seconds)", totalTime)
	}
}
