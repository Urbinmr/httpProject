package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
)

const (
	requestIDKey = iota
	pathKey
	methodKey
)

func beforeIDExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestIDKey, r.Header.Get("X-Request-Id"))
}

func beforePATHExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, pathKey, r.URL.EscapedPath())
}

func beforeMethodExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, methodKey, r.Method)
}

func loggingMiddlware(l log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (result interface{}, err error) {
			var resp string

			defer func(b time.Time) {
				l.Log(
					"timestamp", time.Now(),
					"path", ctx.Value(pathKey),
					"method", ctx.Value(methodKey),
					"result", resp,
					"err", err,
					"request_id", ctx.Value(requestIDKey),
					"elapsed", time.Since(b),
				)
			}(time.Now())

			result, err = next(ctx, request)
			if r, ok := result.(fmt.Stringer); ok {
				resp = r.String()
			}
			return
		}
	}
}
