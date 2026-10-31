// Package logger writes the service's logs as JSON lines on stdout, where Filebeat picks them
// up for Kibana (see deploy/filebeat), and tags each record with the request's transaction id.
//
//   - Level 1 (INFO): one record per HTTP request with its response (Middleware), and one per
//     queue message published or consumed (queue.go, in the services that use RabbitMQ).
//   - Level 2 (DEBUG): also the input and output of the service functions (Trace).
//
// LOG_LEVEL=2 turns level 2 on; otherwise only level 1 is written. This file is copied into
// every service, like the metrics package: keep the copies the same.
package logger

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// Header carries the transaction id. The gateway takes it from the client, or makes one,
	// and each service forwards it on its calls to the others (see Client).
	Header = "X-Transaction-Id"

	// ginKey is where Middleware also stores the id with c.Set: handlers that pass the
	// *gin.Context itself as the context only resolve string keys set that way.
	ginKey = "transaction_id"
)

type ctxKey struct{}

// tracer writes Trace's level 2 records. The default logger stays at level 1 even with
// LOG_LEVEL=2, because third-party packages (the RabbitMQ client) write their debug output to it.
var tracer *slog.Logger

// Setup makes a JSON logger for service the default. That also sends the standard log
// package's output (log.Printf, ...) through it, as level 1 records.
func Setup(service string) {
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "2" {
		level = slog.LevelDebug
	}
	// One JSON handler under both loggers, so their lines never interleave on stdout.
	out := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Filebeat decodes each line over its own "message" field, so the decoded record
			// replaces the raw line instead of being stored next to it.
			if len(groups) == 0 && a.Key == slog.MessageKey {
				a.Key = "message"
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler{out, slog.LevelInfo}).With("service", service))
	tracer = slog.New(handler{out, level}).With("service", service)
}

// handler drops records below level and adds the transaction id of the record's context.
type handler struct {
	slog.Handler
	level slog.Level
}

func (h handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h handler) Handle(ctx context.Context, r slog.Record) error {
	if id := TransactionID(ctx); id != "" {
		r.AddAttrs(slog.String("transaction_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return handler{h.Handler.WithAttrs(attrs), h.level}
}

func (h handler) WithGroup(name string) slog.Handler {
	return handler{h.Handler.WithGroup(name), h.level}
}

// WithTransactionID returns a copy of ctx that carries id.
func WithTransactionID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// TransactionID returns the transaction id ctx carries, or "".
func TransactionID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(ctxKey{}).(string); ok {
		return id
	}
	id, _ := ctx.Value(ginKey).(string)
	return id
}

// NewTransactionID returns a random UUID (version 4).
func NewTransactionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Middleware logs each request with its response at level 1. It takes the transaction id from
// the X-Transaction-Id header, or makes one if there is none, puts it in the request's context
// and sends it back in the response header.
// Register it before gin.Recovery, so a recovered panic is logged with its 500.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		id := c.GetHeader(Header)
		if id == "" {
			id = NewTransactionID()
		}
		c.Set(ginKey, id)
		ctx := WithTransactionID(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(Header, id)

		requestBody := readRequestBody(c.Request)
		writer := &bodyWriter{ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		status := writer.Status()
		level := slog.LevelInfo
		switch {
		case status >= 500:
			level = slog.LevelError
		case status >= 400:
			level = slog.LevelWarn
		}
		slog.Log(ctx, level, fmt.Sprintf("%s %s %d", c.Request.Method, c.Request.URL.Path, status),
			"kind", "http",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"route", c.FullPath(),
			"query", c.Request.URL.RawQuery,
			"status", status,
			"duration_ms", milliseconds(time.Since(start)),
			"client_ip", c.ClientIP(),
			"request_body", requestBody,
			"response_body", formatBody(writer.Header().Get("Content-Type"), writer.body.Bytes()),
		)
	}
}

// readRequestBody returns the request body for the log and puts it back, so the handler can
// still read it.
func readRequestBody(r *http.Request) string {
	if r.Body == nil || r.Body == http.NoBody {
		return ""
	}
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	return formatBody(r.Header.Get("Content-Type"), body)
}

// bodyWriter keeps a copy of the response body for the log.
type bodyWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// Client is for calls to other services: it forwards the transaction id of each request's
// context in the X-Transaction-Id header.
var Client = &http.Client{Transport: forwardTransactionID{http.DefaultTransport}}

type forwardTransactionID struct{ next http.RoundTripper }

func (t forwardTransactionID) RoundTrip(r *http.Request) (*http.Response, error) {
	if id := TransactionID(r.Context()); id != "" && r.Header.Get(Header) == "" {
		r = r.Clone(r.Context()) // a RoundTripper must not change the caller's request
		r.Header.Set(Header, id)
	}
	return t.next.RoundTrip(r)
}

// Trace logs a call of a service function at level 2, with its input and output. Defer it at
// the top of the function with pointers to the named results, so it sees their final values:
//
//	func (s OrderService) GetOrderById(id int, ctx context.Context) (order model.Order, err error) {
//		defer logger.Trace(ctx, "OrderService.GetOrderById", id)(&order, &err)
//
// A last result of type error goes in the record's error field instead of its output.
func Trace(ctx context.Context, function string, input ...any) func(output ...any) {
	if tracer == nil || !tracer.Enabled(ctx, slog.LevelDebug) {
		return func(...any) {}
	}
	start := time.Now()
	// Encoded now: the function may change what its arguments point to.
	in := encode(input)

	return func(output ...any) {
		values := make([]any, len(output))
		for i, p := range output {
			values[i] = deref(p)
		}
		attrs := []any{"kind", "function", "function", function, "function_input", in}
		if n := len(output); n > 0 {
			if _, isErr := output[n-1].(*error); isErr {
				if err, _ := values[n-1].(error); err != nil {
					attrs = append(attrs, "error", err.Error())
				}
				values = values[:n-1]
			}
		}
		attrs = append(attrs, "function_output", encode(values), "duration_ms", milliseconds(time.Since(start)))
		tracer.DebugContext(ctx, function, attrs...)
	}
}

// deref returns what p points to.
func deref(p any) any {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return p
	}
	return v.Elem().Interface()
}

// encode returns values as a JSON array for the log, with errors as their message.
func encode(values []any) string {
	if values == nil {
		values = []any{}
	}
	for i, v := range values {
		if err, ok := v.(error); ok {
			values[i] = err.Error()
		}
	}
	b, err := json.Marshal(values)
	if err != nil {
		return fmt.Sprintf("%+v", values)
	}
	return string(b)
}

// formatBody returns a body as it is logged: text only.
func formatBody(contentType string, b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if contentType != "" && !strings.Contains(contentType, "json") && !strings.Contains(contentType, "text") &&
		!strings.Contains(contentType, "x-www-form-urlencoded") {
		return fmt.Sprintf("(%s body not logged)", contentType)
	}
	return string(b)
}

func milliseconds(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}
