package fetch

import (
	"net/http"
	"testing"

	"github.com/thekhanj/digikala-api/cli/internal"
)

func TestFetchProducts(t *testing.T) {
	dir := "test/.cache/product"
	f, err := NewProducts(
		http.DefaultClient,
		[]string{"https://api.digikala.com/v2/product/17986495/"},
		dir,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	err = internal.AssertAtLeastOneFile(dir)
	if err != nil {
		t.Fatal(err)
	}
}
