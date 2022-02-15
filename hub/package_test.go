package hub

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"

	"github.com/go-logr/stdr"
	"golang.org/x/net/publicsuffix"
)

var testLog = stdr.NewWithOptions(log.New(os.Stderr, "", log.LstdFlags), stdr.Options{})

func newHTTPClient(t *testing.T) *http.Client {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return &http.Client{
		Jar: jar,
	}
}
