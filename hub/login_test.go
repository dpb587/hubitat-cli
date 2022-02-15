package hub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClientLogin_Success(t *testing.T) {
	var handledLogin, handledRedirect int

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/login":
				handledLogin++

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
				} else if _e, _a := "test-username-1", r.FormValue("username"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				} else if _e, _a := "test-password-1", r.FormValue("password"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				} else if _e, _a := "Login", r.FormValue("submit"); _e != _a {
					t.Fatalf("expected `%v` but got: %v", _e, _a)
				}

				http.SetCookie(
					w,
					&http.Cookie{
						Name:  "TESTCOOKIE",
						Value: "test-session-1",
					},
				)

				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
			case "/":
				handledRedirect++

				cookie, err := r.Cookie("TESTCOOKIE")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				} else if _e, _a := "test-session-1", cookie.Value; _e != _a {
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

	err := hubClient.Login(context.Background(), "test-username-1", "test-password-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _e, _a := 1, handledLogin; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := 1, handledRedirect; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func TestClientLogin_BadResponseStatusCode(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/login":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}),
	)

	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	hubClient := NewClient(testLog, newHTTPClient(t), serverURL)

	err := hubClient.Login(context.Background(), "test-username-1", "test-password-1")
	if err == nil {
		t.Fatalf("expected error but got: nil")
	}
}

func TestClientLogin_RequiresJar(t *testing.T) {
	hubClient := NewClient(testLog, http.DefaultClient, nil)

	err := hubClient.Login(context.Background(), "test-username-1", "test-password-1")
	if err == nil {
		t.Fatalf("expected error but got: nil")
	} else if _e, _a := "cookie jar is missing", err.Error(); !strings.Contains(_a, _e) {
		t.Fatalf("expected substring `%s` but got: %v", _e, _a)
	}
}
