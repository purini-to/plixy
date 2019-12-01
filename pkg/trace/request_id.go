package trace

import (
	"context"
	"net/http"
)

type reqIDKeyType int

const (
	RequestIDHeader              = "x-request-id"
	reqIDContextKey reqIDKeyType = iota
)

func RequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, reqIDContextKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(reqIDContextKey).(string); ok {
		return id
	}
	return ""
}

func RequestIDFromRequest(r *http.Request) string {
	return r.Header.Get(RequestIDHeader)
}

func RequestIDToReqRes(w http.ResponseWriter, r *http.Request, requestID string) {
	r.Header.Set(RequestIDHeader, requestID)
	w.Header().Set(RequestIDHeader, requestID)
}
