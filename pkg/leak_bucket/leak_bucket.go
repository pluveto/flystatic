// 使用示例：
// func main() {
// 	// 配置静态文件目录和端口
// 	fs := http.FileServer(http.Dir("static"))
// 	port := "8080"

// 	// 定义限速参数
// 	rate := 1024.0                     // 每秒限制写入1KB
// 	bucketSize := 1024.0               // 桶的容量为1KB
// 	interval := time.Millisecond * 100 // 每次写入间隔100毫秒

// 	// 定义处理函数
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		rateLimitedWriter := &RateLimitedResponseWriter{
// 			Writer:         w,
// 			ResponseWriter: w,
// 			Limit:          rate,
// 			Bucket:         NewLeakyBucket(rate, bucketSize),
// 			Interval:       interval,
// 		}
// 		fs.ServeHTTP(rateLimitedWriter, r)
// 	}

// 	// 创建 HTTP 服务器并开始监听端口
// 	server := &http.Server{
// 		Addr:         ":" + port,
// 		Handler:      http.HandlerFunc(handler),
// 		ReadTimeout:  5 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 	}
// 	fmt.Printf("Server started on port %s\n", port)
// 	server.ListenAndServe()
// }

package leak_bucket

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type LeakyBucket struct {
	Rate       float64   // 每秒速率
	BucketSize float64   // 漏桶的容量
	LastLeakAt time.Time // 上次漏水时间
	Remaining  float64   // 漏桶中的水量
}

// NewLeakyBucket 创建一个漏桶，漏桶的容量为 bucketSize 字节，漏水速率为 rate 字节/秒
func NewLeakyBucket(rate, bucketSize float64) *LeakyBucket {
	return &LeakyBucket{
		Rate:       rate,
		BucketSize: bucketSize,
		LastLeakAt: time.Now(),
	}
}

// Allow 尝试加水，如果加水后水量超过了桶的容量，则返回 false，否则返回 true
func (l *LeakyBucket) Allow(amount float64) bool {
	// 首先漏水
	now := time.Now()
	timeElapsed := now.Sub(l.LastLeakAt).Seconds()
	leakedAmount := timeElapsed * l.Rate
	l.Remaining = l.Remaining - leakedAmount
	if l.Remaining < 0 {
		l.Remaining = 0
	}
	l.LastLeakAt = now

	// 然后尝试加水
	if l.Remaining+amount <= l.BucketSize {
		l.Remaining = l.Remaining + amount
		return true
	} else {
		return false
	}
}

// RateLimitedResponseWriter 限速的 ResponseWriter
type RateLimitedResponseWriter struct {
	io.Writer
	http.ResponseWriter
	Limit    float64
	Bucket   *LeakyBucket
	Interval time.Duration
}

// Write 限速写入
func (w *RateLimitedResponseWriter) Write(p []byte) (int, error) {
	if !w.Bucket.Allow(float64(len(p))) {
		return 0, fmt.Errorf("rate limit exceeded")
	}
	time.Sleep(w.Interval)
	n, err := w.Writer.Write(p)
	return n, err
}
