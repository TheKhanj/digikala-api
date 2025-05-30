package config

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/thekhanj/digikala-api/cli/internal"
)

func (this ProxyCmd) logStderr(bytes []byte) {
	lines := strings.Split(string(bytes), "\n")

	for _, line := range lines {
		log.Printf("proxy-cmd: %s: %s", string(this), line)
	}
}

func (this ProxyCmd) Execute() ([]string, error) {
	arr := strings.Split(string(this)[1:], " ")
	name := arr[0]
	absName := internal.GetAbsPath(name)
	args := make([]string, 0)
	for _, arg := range arr[1:] {
		if strings.TrimSpace(arg) != "" {
			args = append(args, arg)
		}
	}

	cmd := exec.CommandContext(context.Background(), absName, args...)
	bytes, err := cmd.Output()
	if err != nil {
		this.logStderr(bytes)

		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")
	ret := make([]string, 0)
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			ret = append(ret, line)
		}
	}

	return ret, nil
}

func (this *ConfigApiClient) GetProxies() ([]string, error) {
	ret := make([]string, 0)

	for _, proxy := range this.Proxies {
		p := proxy.(string)

		if v := Proxy(p); v.UnmarshalJSON([]byte("\""+v+"\"")) == nil {
			ret = append(ret, proxy.(string))
		} else if v := ProxyCmd(p); v.UnmarshalJSON([]byte("\""+v+"\"")) == nil {
			arr, err := v.Execute()
			if err != nil {
				return nil, err
			}
			ret = append(ret, arr...)
		} else {
			return nil, fmt.Errorf("invalid proxy type: %v", proxy)
		}

	}

	return ret, nil
}
