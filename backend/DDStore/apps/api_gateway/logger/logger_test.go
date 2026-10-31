package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// captureLogs makes both loggers write JSON records, level 2 included, into the returned buffer.
func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	oldDefault, oldTracer := slog.Default(), tracer
	l := slog.New(handler{slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}), slog.LevelDebug})
	slog.SetDefault(l)
	tracer = l
	t.Cleanup(func() { slog.SetDefault(oldDefault); tracer = oldTracer })
	return &buf
}

func TestSetupLevel2OnlyForTrace(t *testing.T) {
	oldDefault, oldTracer := slog.Default(), tracer
	t.Cleanup(func() { slog.SetDefault(oldDefault); tracer = oldTracer })
	t.Setenv("LOG_LEVEL", "2")
	Setup("test")

	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		t.Error("the default logger, which libraries use too, should stay at level 1")
	}
	if !tracer.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("Trace should log at level 2 with LOG_LEVEL=2")
	}
}

func records(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var r map[string]any
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			t.Fatalf("log line is not JSON: %q", line)
		}
		out = append(out, r)
	}
	return out
}

func TestMiddlewareKeepsTransactionIdAndBody(t *testing.T) {
	buf := captureLogs(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware())
	var seenBody, seenID string
	r.POST("/login", func(c *gin.Context) {
		b, _ := io.ReadAll(c.Request.Body)
		seenBody, seenID = string(b), TransactionID(c)
		c.JSON(http.StatusOK, gin.H{"access_token": "abc", "ok": true})
	})

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"name":"bob","password":"hunter2"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(Header, "tx-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if seenBody != `{"name":"bob","password":"hunter2"}` {
		t.Errorf("handler got body %q, want it unchanged", seenBody)
	}
	if seenID != "tx-123" || w.Header().Get(Header) != "tx-123" {
		t.Errorf("transaction id: handler %q, response header %q, want tx-123", seenID, w.Header().Get(Header))
	}

	rec := records(t, buf)[0]
	if rec["transaction_id"] != "tx-123" || rec["status"] != float64(200) || rec["kind"] != "http" {
		t.Errorf("unexpected record %v", rec)
	}
	if body := rec["request_body"]; body != `{"name":"bob","password":"hunter2"}` {
		t.Errorf("request body logged as %v, want it as sent", body)
	}
	if body := rec["response_body"]; body != `{"access_token":"abc","ok":true}` {
		t.Errorf("response body logged as %v, want it as sent", body)
	}
}

func TestMiddlewareTransactionId(t *testing.T) {
	captureLogs(t)
	r := gin.New()
	r.Use(Middleware())
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	send := func(id string) string {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if id != "" {
			req.Header.Set(Header, id)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Header().Get(Header)
	}

	if got := send("any value, kept as sent"); got != "any value, kept as sent" {
		t.Errorf("got transaction id %q, want the one sent", got)
	}
	if got := send(""); len(got) != 36 {
		t.Errorf("without a header: got transaction id %q, want a new UUID", got)
	}
}

func TestClientForwardsTransactionId(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { got = r.Header.Get(Header) }))
	defer srv.Close()

	req, _ := http.NewRequestWithContext(WithTransactionID(context.Background(), "tx-9"), http.MethodGet, srv.URL, nil)
	resp, err := Client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got != "tx-9" {
		t.Errorf("service got transaction id %q, want tx-9", got)
	}
}

func TestTraceLogsInputOutputAndError(t *testing.T) {
	buf := captureLogs(t)
	ctx := WithTransactionID(context.Background(), "tx-1")
	f := func(id int) (name string, err error) {
		defer Trace(ctx, "Svc.Get", id)(&name, &err)
		return "", errors.New("not found")
	}
	f(7)

	rec := records(t, buf)[0]
	if rec["function"] != "Svc.Get" || rec["function_input"] != "[7]" || rec["function_output"] != `[""]` ||
		rec["error"] != "not found" || rec["transaction_id"] != "tx-1" || rec["level"] != "DEBUG" {
		t.Errorf("unexpected record %v", rec)
	}
}
