package http

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	client := resty.New()

	for i := 0; i < b.N; i++ {
		resp, err := client.GetClient().Get("http://www.baidu.com")
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
		resp, err := http.Get("http://www.baidu.com")
		if err != nil {
			b.Error(err)
		}

		if resp.StatusCode != 200 {
			b.Error("status code is not 200")
		}
	}
}

func BenchmarkFastHttp(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}
