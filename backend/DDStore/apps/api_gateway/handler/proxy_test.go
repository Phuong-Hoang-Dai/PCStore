package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/logger"
	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

// fakeService points *service at a test server running handler until the test ends.
func fakeService(t *testing.T, service **url.URL, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	old := *service
	*service, _ = url.Parse(srv.URL)
	t.Cleanup(func() { srv.Close(); *service = old })
}

func TestForwardToPassesRequestAndResponseThrough(t *testing.T) {
	var got *http.Request
	var gotBody string
	fakeService(t, &productService, func(w http.ResponseWriter, r *http.Request) {
		got = r
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set(logger.Header, "tx-7")
		w.Header().Set("X-Custom", "kept")
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(`[1,2,3]`))
	})
	r := gin.New()
	r.Use(logger.Middleware())
	r.PUT("/api/v1/product/:id", forwardTo(productService))

	req := httptest.NewRequest(http.MethodPut, "/api/v1/product/42?x=1&y=2", strings.NewReader(`{"name":"pc"}`))
	req.Header.Set("Cookie", "access_token=abc")
	req.Header.Set(logger.Header, "tx-7")
	w := newRecorder()
	r.ServeHTTP(w, req)

	if got.Method != http.MethodPut || got.URL.Path != "/product/42" || got.URL.RawQuery != "x=1&y=2" || gotBody != `{"name":"pc"}` {
		t.Errorf("service got %s %s?%s %q", got.Method, got.URL.Path, got.URL.RawQuery, gotBody)
	}
	if got.Header.Get("Cookie") != "" || got.Header.Get(logger.Header) != "tx-7" || got.Header.Get("Content-Type") != "application/json" {
		t.Errorf("service got headers %v", got.Header)
	}
	if w.Code != http.StatusTeapot || w.Body.String() != `[1,2,3]` || w.Header().Get("X-Custom") != "kept" {
		t.Errorf("client got %d %q %v", w.Code, w.Body.String(), w.Header())
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "" || len(w.Header().Values(logger.Header)) != 1 {
		t.Errorf("service's CORS and transaction id headers should be dropped: %v", w.Header())
	}
}

func TestForwardSendsGeneratedTransactionId(t *testing.T) {
	var got string
	fakeService(t, &productService, func(w http.ResponseWriter, r *http.Request) { got = r.Header.Get(logger.Header) })
	r := gin.New()
	r.Use(logger.Middleware())
	r.GET("/api/v1/category", forwardTo(productService))

	w := newRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/category", nil))
	if got == "" || got != w.Header().Get(logger.Header) {
		t.Errorf("service got transaction id %q, client got %q: want the same generated id", got, w.Header().Get(logger.Header))
	}
}

func TestForwardListDefaultsPaging(t *testing.T) {
	var query string
	fakeService(t, &productService, func(w http.ResponseWriter, r *http.Request) { query = r.URL.RawQuery })
	r := gin.New()
	r.GET("/api/v1/product", forwardList(productService))

	for sent, want := range map[string]string{"": "limit=0&offset=0", "?limit=5&q=pc": "limit=5&offset=0&q=pc"} {
		r.ServeHTTP(newRecorder(), httptest.NewRequest(http.MethodGet, "/api/v1/product"+sent, nil))
		if query != want {
			t.Errorf("sent %q: service got query %q, want %q", sent, query, want)
		}
	}
}

func TestForwardToUnreachableServiceIs502(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	down, _ := url.Parse(srv.URL)
	srv.Close()
	r := gin.New()
	r.GET("/api/v1/category", forwardTo(down))

	w := newRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/category", nil))
	if w.Code != http.StatusBadGateway || !strings.Contains(w.Body.String(), `"success":false`) {
		t.Errorf("got %d %s, want a 502 error", w.Code, w.Body.String())
	}
}

func TestGetOrderByIdOnlyOwnOrderForUsers(t *testing.T) {
	fakeService(t, &orderService, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"data":{"id":5,"userId":7}}`)
	})

	for userId, want := range map[int]int{7: http.StatusOK, 8: http.StatusForbidden} {
		r := gin.New()
		r.GET("/api/v1/order/:id", func(c *gin.Context) { c.Set("role", "user"); c.Set("id", userId) }, getOrderById())
		w := newRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/order/5", nil))
		if w.Code != want {
			t.Errorf("user %d reading order of user 7: got %d, want %d", userId, w.Code, want)
		}
	}
}

// recorder adds CloseNotify, which ReverseProxy calls through gin's writer: a real connection
// has it, httptest.ResponseRecorder doesn't.
type recorder struct{ *httptest.ResponseRecorder }

func (recorder) CloseNotify() <-chan bool { return make(chan bool) }

func newRecorder() recorder { return recorder{httptest.NewRecorder()} }
