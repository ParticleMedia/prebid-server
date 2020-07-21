// mock_rubicon.go
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// mock ad engine that sleeps for a random amount of time, then
// responds with the JSON string.

var Port = flag.Int("listen-port", 80, "Web server port")

var payload = []byte(``)

// sleep interval is approximately normally distributed:
// 40ms - 200ms
// == 160ms

// mean is 120 (ms)
// P(120-3*s <= x <= 120+3*s) == 0.9973
// 120-3*s == 40, 120+3*s == 200
// 3*s = 120-40 = 80, s = 26.666
func handler(w http.ResponseWriter, r *http.Request) {
	t := time.Duration(1000*(rand.NormFloat64()*26.6666+120.0)) * time.Microsecond
	time.Sleep(t)
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(payload)
}

func main() {
	flag.Parse()

	http.HandleFunc("/openrtb2/auction", handler)
	listen_addr := fmt.Sprintf(":%d", *Port)
	http.ListenAndServe(listen_addr, nil)
}
