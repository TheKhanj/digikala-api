package proxy

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

type DialFunc func(network, addr string) (net.Conn, error)

func socks5Transport(proxyStr string) (*http.Transport, error) {
	p := strings.ReplaceAll(proxyStr, "socks5://", "")

	dialer, err := proxy.SOCKS5("tcp", p, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	return &http.Transport{Dial: dialer.Dial}, nil
}

func httpTransport(proxyStr string) (*http.Transport, error) {
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}

	var tlsCfg *tls.Config = nil
	if strings.HasPrefix(proxyStr, "https://") {
		tlsCfg = &tls.Config{}
	}

	return &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: tlsCfg,
	}, nil
}

// Supported proxy protocols are:
//
//	http://...
//	https://...
//	socks5://...
func NewProxyClient(proxyStr string) (*http.Client, error) {
	var t *http.Transport
	var err error
	if strings.HasPrefix(proxyStr, "socks5://") {
		t, err = socks5Transport(proxyStr)
	} else if strings.HasPrefix(proxyStr, "http://") ||
		strings.HasPrefix(proxyStr, "https://") {
		t, err = httpTransport(proxyStr)
	} else {
		err = fmt.Errorf("not supported proxy protocol: (%s)", proxyStr)
	}

	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: t}, nil
}

func NewProxyClientList(proxies []string) ([]*http.Client, error) {
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
