package util

import (
	"context"
	"net/http"
)

const Retry = "retry"

func SetRetry(r *http.Request, times int) *http.Request {
	ctx := context.WithValue(r.Context(), Retry, times)
	return r.WithContext(ctx)
}

func GetRetry(r *http.Request) int {
	retry, ok := r.Context().Value(Retry).(int)
	if !ok {
		return 0
	}
	return retry
}
