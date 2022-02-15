package hub

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClientBackup_Success(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/backupDB":
				if _e, _a := http.MethodGet, r.Method; _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				} else if _e, _a := "test-file-1", r.URL.Query().Get("fileName"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				w.Header().Set("Content-Disposition", `attachment; filename="actual-name.lzf"`)
				w.Header().Set("Content-Length", "12")
				w.Write([]byte("hello world\n"))
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	meta, reader, err := hubClient.DownloadBackupFile(context.Background(), "test-file-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _e, _a := "actual-name.lzf", meta.Name; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := 12, meta.Size; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}

	readerBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if _e, _a := "hello world\n", string(readerBytes); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func TestClientBackup_BadResponseStatusCode(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/backupDB":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	_, _, err := hubClient.DownloadBackupFile(context.Background(), "test-file-1")
	if err == nil {
		t.Fatalf("expected error but got: nil")
	}
}
