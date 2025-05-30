package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/thekhanj/digikala-api/cli/internal"
	"github.com/thekhanj/digikala-api/cli/proxy"
)

type Products struct {
	client proxy.HttpClient
	urls   []string
	dir    string
}

func (this *Products) fetchProduct(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected http status code: %v", res.StatusCode,
		)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (this *Products) saveBody(index int, body []byte) error {
	j := make(map[string]interface{})
	err := json.Unmarshal(body, &j)
	if err != nil {
		return err
	}

	// TODO: maybe use jq here...
	// FUCK IT USE ID FROM URL
	id := int(
		j["data"].(map[string]interface{})["product"].(map[string]interface{})["id"].(float64),
	)

	fileName := fmt.Sprintf("%05d-%d.json", index, id)

	return os.WriteFile(
		path.Join(this.dir, fileName),
		body, 0644,
	)
}

func (this *Products) Fetch() error {
	for index, url := range this.urls {
		body, err := this.fetchProduct(url)
		if err != nil {
			return err
		}

		err = this.saveBody(index+1, body)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewProducts(
	client proxy.HttpClient,
	urls []string,
	dir string,
) (*Products, error) {
	absDir := internal.GetAbsPath(dir)
	err := os.MkdirAll(absDir, 0755)
	if err != nil {
		return nil, err
	}

	return &Products{client: client, urls: urls, dir: absDir}, nil
}
