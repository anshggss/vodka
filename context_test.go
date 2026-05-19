package vodka

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAbortStopsChain(t *testing.T) {
	calls := []string{}

	// First handler aborts
	handler1 := func(c *Context) {
		calls = append(calls, "h1")
		c.Abort()
	}

	// This should NOT run
	handler2 := func(c *Context) {
		calls = append(calls, "h2")
	}

	c := &Context{handlers: []HandlerFunc{handler1, handler2}, index: -1}
	c.Next()

	if len(calls) != 1 || calls[0] != "h1" {
		t.Errorf("got %v, want [h1]", calls)
	}
}

func TestSetCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	c := &Context{Writer: rr, Request: req}

	c.SetCookie("token", "abc123", 3600)

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Name != "token" || cookies[0].Value != "abc123" {
		t.Errorf("got %s=%s, want token=abc123", cookies[0].Name, cookies[0].Value)
	}
	if cookies[0].MaxAge != 3600 {
		t.Errorf("got MaxAge %d, want 3600", cookies[0].MaxAge)
	}
}

func TestCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "xyz"})
	c := &Context{Writer: rr, Request: req}

	val, err := c.Cookie("session")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "xyz" {
		t.Errorf("got %q, want xyz", val)
	}

	_, err = c.Cookie("missing")
	if err == nil {
		t.Error("expected error for missing cookie")
	}
}

func TestClearCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	c := &Context{Writer: rr, Request: req}

	c.ClearCookie("token")

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].MaxAge != -1 {
		t.Errorf("got MaxAge %d, want -1", cookies[0].MaxAge)
	}
}