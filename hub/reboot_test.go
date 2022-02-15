package hub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClientReboot_Success(t *testing.T) {
	var handledReboot int

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/reboot":
				handledReboot++

				if _e, _a := http.MethodPost, r.Method; _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				w.WriteHeader(http.StatusOK)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.Reboot(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _e, _a := 1, handledReboot; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func TestClientReboot_BadResponseStatusCode(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/reboot":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.Reboot(context.Background())
	if err == nil {
		t.Fatalf("expected error but got: nil")
	}
}
