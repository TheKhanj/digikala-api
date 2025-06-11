package jq

import (
	"strings"
	"testing"
)

func TestJq(t *testing.T) {
	str := `{"name":"pooyan","lastname":"khanjankhani"}`
	filter := `{fullname: (.name + " " + .lastname)}`

	jq := NewJq([]byte(str), filter, "-c")

	res, err := jq.Start()
	if err != nil {
		t.Fatal(err)
	}

	expectedRes := `{"fullname":"pooyan khanjankhani"}`
	if strings.TrimSpace(string(res)) != expectedRes {
		t.Fatalf("unexpected result: %s", string(res))
	}
}
