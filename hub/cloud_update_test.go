package hub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCloudCheckForUpdate_SuccessAvailable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/cloud/checkForUpdate":
				if _e, _a := http.MethodGet, r.Method; _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"version":"2.3.1.142","upgrade":true,"releaseNotesContent":"<div class='post' itemprop='articleBody'>...</div>","status":"UPDATE_AVAILABLE"}`))
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	stat, err := hubClient.CloudCheckForUpdate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _e, _a := "2.3.1.142", stat.Version; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := true, stat.Upgrade; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "<div class='post' itemprop='articleBody'>...</div>", stat.ReleaseNotesContent; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "UPDATE_AVAILABLE", stat.Status; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func TestCloudCheckForUpdate_BadResponseStatusCode(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/cloud/checkForUpdate":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	_, err := hubClient.CloudCheckForUpdate(context.Background())
	if err == nil {
		t.Fatalf("expected error but got: nil")
	}
}

func TestCloudCheckUpdateStatus_SuccessAvailable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/cloud/checkUpdateStatus":
				if _e, _a := http.MethodGet, r.Method; _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"status":"DOWNLOAD_IN_PROGRESS","percent":10}`))
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	stat, err := hubClient.CloudCheckUpdateStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _e, _a := "DOWNLOAD_IN_PROGRESS", stat.Status; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := 10, stat.Percent; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func TestCloudCheckUpdateStatus_BadResponseStatusCode(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/cloud/checkUpdateStatus":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	_, err := hubClient.CloudCheckUpdateStatus(context.Background())
	if err == nil {
		t.Fatalf("expected error but got: nil")
	}
}

func TestCloudUpdatePlatform_Success(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/cloud/updatePlatform":
				if _e, _a := http.MethodGet, r.Method; _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"success":"true"}`))
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.CloudUpdatePlatform(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloudUpdatePlatform_NotSuccess(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/cloud/updatePlatform":
				if _e, _a := http.MethodGet, r.Method; _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"success":"false"}`))
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.CloudUpdatePlatform(context.Background())
	if err == nil {
		t.Fatalf("expected error but got: nil")
	} else if _e, _a := "update request failed", err.Error(); !strings.Contains(_a, _e) {
		t.Fatalf("expected substring `%v` but got: %v", _e, _a)
	}
}

func TestCloudUpdatePlatform_BadResponseStatusCode(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/cloud/updatePlatform":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.CloudUpdatePlatform(context.Background())
	if err == nil {
		t.Fatalf("expected error but got: nil")
	}
}
