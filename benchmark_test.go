package vodka

import (
	"net/http/httptest"
	"testing"
)

func BenchmarkRouterStatic(b *testing.B) {
	app := NewRouter()

	app.GET("/test", func(c *Context) {
		c.String(200, "hello")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

func BenchmarkRoutingOnly(b *testing.B) {
	app := NewRouter()

	app.GET("/users/:id", func(c *Context) {})

	req := httptest.NewRequest("GET", "/users/123", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

func BenchmarkDeepParam(b *testing.B) {
	app := NewRouter()

	app.GET("/api/v1/projects/:project/user/:userId/objective/:obj/comment/:comment", func(c *Context) {
		_ = c.Param("project")
		_ = c.Param("userId")
		_ = c.Param("obj")
		_ = c.Param("comment")
	})

	req := httptest.NewRequest("GET", "/api/v1/projects/vodka/user/123/objective/test/comment/67", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}
