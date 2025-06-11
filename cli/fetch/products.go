package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/thekhanj/digikala-api/cli/internal"
	"github.com/thekhanj/digikala-api/cli/jq"
	"github.com/thekhanj/digikala-api/cli/proxy"
)

type Products struct {
	client proxy.HttpClient
	urls   []string
	dir    string
}

func (this *Products) fetchProduct(url string) ([]byte, error) {
	log.Printf("products: fetching %s", url)

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

	return jq.NewJq(bytes, ".").Start()
}

func (this *Products) saveBody(index int, url string, body []byte) error {
	j := make(map[string]interface{})
	err := json.Unmarshal(body, &j)
	if err != nil {
		return err
	}

	id, err := this.getIdFromUrl(url)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%05d-%d.json", index, id)

	return os.WriteFile(
		path.Join(this.dir, fileName),
		body, 0644,
	)
}

func (this *Products) getIdFromUrl(url string) (int, error) {
	r := regexp.MustCompile("/v2/product/(?P<id>[0-9]*)")
	matches := r.FindStringSubmatch(url)

	subs := r.SubexpNames()
	idFound := false
	var idStr string

	for i, sub := range subs {
		if sub == "id" {
			idFound = true
			idStr = matches[i]
			break
		}
	}

	if !idFound {
		return 0, errors.New("missing product id in url")
	}

	return strconv.Atoi(idStr)
}

func (this *Products) Fetch() error {
	for index, url := range this.urls {
		body, err := this.fetchProduct(url)
		if err != nil {
			return err
		}

		err = this.saveBody(index+1, url, body)
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
