package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

func (this *ConfigApiFetch) getIdFromUrl(url string) (int, error) {
	r := regexp.MustCompile("^.*/product/dkp-(?P<id>[0-9]*).*$")
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

func (this *ConfigApiFetch) GetProductsApiUrls() ([]string, error) {
	ret := make([]string, 0)
	for _, url := range this.Products {
		id, err := this.getIdFromUrl(url)
		if err != nil {
			return nil, err
		}

		ret = append(
			ret,
			fmt.Sprintf(
				"https://api.digikala.com/v2/product/%d/", id,
			),
		)
	}

	return ret, nil
}
