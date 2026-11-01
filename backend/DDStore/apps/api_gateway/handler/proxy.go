package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/logger"
	"github.com/gin-gonic/gin"
)

// The services the gateway forwards to, by their compose service names.
var (
	userService    = &url.URL{Scheme: "http", Host: "user_service:6666"}
	productService = &url.URL{Scheme: "http", Host: "product_service:8888"}
	orderService   = &url.URL{Scheme: "http", Host: "order_service:7777"}
)

// forwardTo returns a handler that forwards the request to the same path on service, without
// /api/v1: GET /api/v1/product/42?x=1 goes to GET {service}/product/42?x=1.
func forwardTo(service *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		forward(c, service, strings.TrimPrefix(c.Request.URL.Path, "/api/v1"), c.Request.URL.RawQuery)
	}
}

// forwardList is forwardTo for list routes; see pagingQuery.
func forwardList(service *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		forward(c, service, strings.TrimPrefix(c.Request.URL.Path, "/api/v1"), pagingQuery(c))
	}
}

// pagingQuery returns the request's query with a missing limit or offset set to 0, as the
// gateway always sent them; catalog_service would otherwise default limit to 10.
func pagingQuery(c *gin.Context) string {
	query := c.Request.URL.Query()
	for _, key := range []string{"limit", "offset"} {
		if query.Get(key) == "" {
			query.Set(key, "0")
		}
	}
	return query.Encode()
}

// forward sends the request to path on service with query as its query string, and its
// method, headers and body as they are, then streams the service's response back unchanged.
// If the service can't be reached, the client gets a 502.
func forward(c *gin.Context, service *url.URL, path, query string) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(service)
			r.Out.URL.Path, r.Out.URL.RawPath, r.Out.URL.RawQuery = path, "", query
			// The service then logs the client's address instead of the gateway's.
			r.SetXForwarded()
			// The services don't use the login cookies, and they read bodies as JSON.
			r.Out.Header.Del("Cookie")
			if r.Out.Header.Get("Content-Type") == "" {
				r.Out.Header.Set("Content-Type", "application/json")
			}
			// An id the client sent is already in the copied headers; one the gateway made is
			// only in the request's context.
			if r.Out.Header.Get(logger.Header) == "" {
				r.Out.Header.Set(logger.Header, logger.TransactionID(r.In.Context()))
			}
		},
		ModifyResponse: dropGatewayHeaders,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			responeError(http.StatusBadGateway, err, c)
		},
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

// dropGatewayHeaders removes from a service's response the headers the gateway has already
// set: CORS, for the frontend's origin (catalog_service also answers "*", and two values make
// the browser reject the response), and the transaction id.
func dropGatewayHeaders(res *http.Response) error {
	for name := range res.Header {
		if strings.HasPrefix(name, "Access-Control-") {
			res.Header.Del(name)
		}
	}
	res.Header.Del(logger.Header)
	return nil
}

// fetch calls a service and returns its status and body, for the routes that need to read
// the response before answering. logger.Client forwards the transaction id of ctx.
func fetch(ctx context.Context, method, url string, body io.Reader) (status int, data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := logger.Client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	return resp.StatusCode, data, err
}

func responeError(errCode int, err error, ctx *gin.Context) {
	ctx.JSON(errCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
