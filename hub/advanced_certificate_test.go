package hub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClientUpdateAdvancedCertificate_Success(t *testing.T) {
	var handledSave, handledRedirect int

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/advanced/certificate/save":
				handledSave++

				if _e, _a := http.MethodPost, r.Method; _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				} else if _e, _a := "application/x-www-form-urlencoded", r.Header.Get("Content-Type"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				err := r.ParseForm()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				} else if _e, _a := 3, len(r.Form); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				} else if _e, _a := "test-certificate-1", r.FormValue("certificate"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				} else if _e, _a := "test-privateKey-1", r.FormValue("privateKey"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				} else if _e, _a := "Save Certificate and Key", r.FormValue("_action_update"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
			case "/":
				handledRedirect++

				w.WriteHeader(http.StatusOK)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.UpdateAdvancedCertificates(context.Background(), []byte("test-certificate-1"), []byte("test-privateKey-1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _e, _a := 1, handledSave; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := 1, handledRedirect; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func TestClientUpdateAdvancedCertificate_BadResponseStatusCode(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/hub/advanced/certificate/save":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.UpdateAdvancedCertificates(context.Background(), []byte("test-certificate-1"), []byte("test-privateKey-1"))
	if err == nil {
		t.Fatalf("expected error but got: nil")
	}
}
