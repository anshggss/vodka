package vodka

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestRecovery_NoPanicWritesNormally(t *testing.T) {
    engine := NewRouter()
    engine.Use(Logger(), Recovery(), ErrorHandler())

    engine.GET("/ok", func(c *Context) {
        c.JSON(http.StatusOK, M{"status": "ok"})
    })

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/ok", nil)
    engine.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}

func TestRecovery_PanicBeforeWrite(t *testing.T) {
	engine := NewRouter()
	engine.Use(Logger(), Recovery())

	engine.GET("/panic", func(c *Context) {
		panic("something went wrong")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Internal Server Error") {
		t.Errorf("expected error body, got: %s", body)
	}
}

func TestRecovery_PanicAfterWrite(t *testing.T) {
    engine := NewRouter()
    engine.Use(Logger(), Recovery(), ErrorHandler())

    engine.GET("/panic-after-write", func(c *Context) {
        c.JSON(http.StatusOK, M{"data": "partial"})
        panic("panic after write")
    })

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/panic-after-write", nil)
    engine.ServeHTTP(w, req)

    // Headers already sent, Recovery should NOT append a second body
    body := w.Body.String()
    if strings.Count(body, "{") > 1 {
        t.Errorf("expected single JSON body, got multiple concatenated: %s", body)
    }
}