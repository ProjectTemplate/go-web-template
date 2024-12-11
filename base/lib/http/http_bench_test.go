package http

import (
	"github.com/go-resty/resty/v2"
	"go-web-template/base/lib/config"
	"net/http"
	"testing"
	"time"
)

func BenchmarkGet(b *testing.B) {
	client := resty.New()

	for i := 0; i < b.N; i++ {
		resp, err := client.GetClient().Get("http://127.0.0.1:8080/ping")
		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != 200 {
			b.Error("status code is not 200")
		}
	}
}

func BenchmarkHttp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://127.0.0.1:8080/ping")
		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != 200 {
			b.Error("status code is not 200")
		}
	}
}

func BenchmarkFastHttp(b *testing.B) {
	Init(config.FastHttp{
		ReadTimeOut:         time.Second * 10,
		WriteTimeOut:        time.Second * 10,
		MaxIdleConnDuration: time.Hour,
		MaxConnsPerHost:     200,
	})

	for i := 0; i < b.N; i++ {
		err := simpleGet("http://127.0.0.1:8080/ping")
		if err != nil {
			b.Error(err)
		}
	}
}
