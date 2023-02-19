package middleware

import (
	"io"
	"math"
	"net/http"
	"time"

	"github.com/pluveto/flystatic/pkg/leak_bucket"
)

type RateLimitedResponseWriter struct {
	io.Writer
	http.ResponseWriter
	Limit    float64
	Bucket   *leak_bucket.LeakyBucket
	Interval time.Duration
}

func (w *RateLimitedResponseWriter) Write(p []byte) (int, error) {
	if !w.Bucket.Allow(float64(len(p))) {
		return 0, nil
	}
	time.Sleep(w.Interval)
	n, err := w.Writer.Write(p)
	return n, err
}

func NewSpeedLimiter(limit float64, capacity float64) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if math.Abs(limit) < 1 {
			return
		}

		bucket := leak_bucket.NewLeakyBucket(limit, capacity)
		w = &RateLimitedResponseWriter{
			Writer:         w,
			ResponseWriter: w,
			Limit:          limit,
			Bucket:         bucket,
			Interval:       time.Millisecond * 100,
		}
	}
}
