package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/google/go-querystring/query"
	"github.com/valyala/fasthttp"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var headerContentTypeJson = []byte("application/json")

var client *fasthttp.Client
var clientOnce sync.Once

type Entity struct {
	Name string
	Id   int
}

// Init 初始化客户端配置
func Init(config config.FastHttp) {
	clientOnce.Do(func() {
		client = &fasthttp.Client{
			ReadTimeout:                   config.ReadTimeOut,
			WriteTimeout:                  config.ReadTimeOut,
			MaxIdleConnDuration:           config.MaxIdleConnDuration,
			NoDefaultUserAgentHeader:      false, // default User-Agent: fasthttp
			DisableHeaderNamesNormalizing: true,  // If you set the case on your headers correctly you can enable this
			DisablePathNormalizing:        true,
			// increase DNS cache time to an hour instead of default minute
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: time.Minute * 10,
			}).Dial,
		}
	})
}

// Get Get请求
//
// params: 请求参数，为结构体类型需要在字段后面加 url tag
//
//	type request struct {
//		Name    string   `url:"name"`
//		Age     int      `url:"age"`
//		Friends []string `url:"friends"`
//	}
func Get(ctx context.Context, requestUrl string, params interface{}, headers map[string]string, result interface{}) error {
	ctx = utils.WithChildSpan(ctx, "get:"+requestUrl)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	//解析验证url
	_, err := url.Parse(requestUrl)
	if err != nil {
		logger.SpanFailed(ctx, "parse url failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(err))
		return err
	}

	//将结构体转换为url.Values，需要在结构体中添加tag:url
	queryValues := url.Values{}
	if params != nil {
		queryValues, err = query.Values(params)
		if err != nil {
			logger.SpanFailed(ctx, "parse query params failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(err))
			return err
		}
	}

	//设置参数

	logger.Info(ctx, "request url", zap.String("url", requestUrl), zap.Any("header", headers))

	req.SetRequestURI(fmt.Sprintf("%s?%s", requestUrl, queryValues.Encode()))
	req.Header.SetMethod(fasthttp.MethodGet)

	//请求数据
	err = client.Do(req, resp)
	if err != nil {
		logger.SpanFailed(ctx, "http get failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(err))
		return err
	}

	//状态码验证
	statusCode := resp.StatusCode()
	if statusCode != http.StatusOK {
		errInner := errors.New("data request failed , code:" + strconv.Itoa(statusCode))
		logger.SpanFailed(ctx, "http get failed", zap.Int("code", statusCode), zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(errInner))
		return errInner
	}

	//反序列化参数
	err = sonic.Unmarshal(resp.Body(), result)
	if err != nil {
		logger.SpanFailed(ctx, "json unmarshal failed", zap.Int("code", statusCode), zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.String("data", string(resp.Body())), zap.Error(err))
		return err
	}

	logger.SpanSuccess(ctx, "http get success", zap.Int("code", statusCode), zap.String("url", requestUrl), zap.Any("header", headers))

	return nil
}

func sendPostRequest() {
	// per-request timeout
	reqTimeout := time.Duration(100) * time.Millisecond

	reqEntity := &Entity{
		Name: "New entity",
	}
	reqEntityBytes, _ := json.Marshal(reqEntity)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost:8080/")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(reqEntityBytes)

	resp := fasthttp.AcquireResponse()
	err := client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		errName, known := httpConnError(err)
		if known {
			fmt.Fprintf(os.Stderr, "WARN conn error: %v\n", errName)
		} else {
			fmt.Fprintf(os.Stderr, "ERR conn failure: %v %v\n", errName, err)
		}

		return
	}

	statusCode := resp.StatusCode()
	respBody := resp.Body()
	fmt.Printf("DEBUG Response: %s\n", respBody)

	if statusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "ERR invalid HTTP response code: %d\n", statusCode)

		return
	}

	respEntity := &Entity{}
	err = json.Unmarshal(respBody, respEntity)
	if err == nil || errors.Is(err, io.EOF) {
		fmt.Printf("DEBUG Parsed Response: %v\n", respEntity)
	} else {
		fmt.Fprintf(os.Stderr, "ERR failed to parse response: %v\n", err)
	}
}

func httpConnError(err error) (string, bool) {
	var (
		errName string
		known   = true
	)

	switch {
	case errors.Is(err, fasthttp.ErrTimeout):
		errName = "timeout"
	case errors.Is(err, fasthttp.ErrNoFreeConns):
		errName = "conn_limit"
	case errors.Is(err, fasthttp.ErrConnectionClosed):
		errName = "conn_close"
	case reflect.TypeOf(err).String() == "*net.OpError":
		errName = "timeout"
	default:
		known = false
	}

	return errName, known
}
