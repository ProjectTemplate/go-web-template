package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	"github.com/xiaotianfork/go-querystring-json/query"
	"go.uber.org/zap"

	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

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
			MaxIdemponentCallAttempts:     config.RetryTimes,
			ReadTimeout:                   config.ReadTimeOut,
			WriteTimeout:                  config.ReadTimeOut,
			MaxIdleConnDuration:           config.MaxIdleConnDuration,
			MaxConnsPerHost:               config.MaxConnsPerHost,
			NoDefaultUserAgentHeader:      false, // default User-Agent: fasthttp
			DisableHeaderNamesNormalizing: true,  // If you set the case on your headers correctly you can enable this
			DisablePathNormalizing:        true,
			RetryIfErr: func(request *fasthttp.Request, attempts int, err error) (resetTimeout bool, retry bool) {
				//幂等方法
				methodNeedRetry := request.Header.IsGet() || request.Header.IsHead() || request.Header.IsPut()
				if methodNeedRetry {
					return true, true
				}
				return false, false
			},
			// increase DNS cache time to an hour instead of default minute
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: time.Minute * 10,
			}).Dial,
		}
	})
}

func simpleGet(url string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	//请求数据
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	return client.Do(req, resp)
}

// Get Get请求，不需要设置超时时间
//
// params: 请求参数，为结构体类型需要在字段后面加 json tag，会自动转换为对应的参数
//
//	type request struct {
//		Name    string   `json:"name"`
//		Age     int      `json:"age"`
//		Friends []string `json:"friends"`
//	}
func Get(ctx context.Context, requestUrl string, params interface{}, headers map[string]string, result interface{}) error {
	return GetTimeOut(ctx, requestUrl, params, headers, time.Duration(0), result)
}

// Post Post请求，不需要设置超时时间
//
// params: 请求参数，为结构体类型需要在字段后面加 json tag
//
//	type request struct {
//		Name    string   `json:"name"`
//		Age     int      `json:"age"`
//		Friends []string `json:"friends"`
//	}
func Post(ctx context.Context, requestUrl string, params interface{}, body interface{}, headers map[string]string, result interface{}) error {
	return PostTimeOut(ctx, requestUrl, params, body, headers, time.Duration(0), result)
}

// GetTimeOut Get请求，需要设置超时时间
//
// params: 请求参数，为结构体类型需要在字段后面加 json tag，会自动转换为对应的参数
//
//	type request struct {
//		Name    string   `json:"name"`
//		Age     int      `json:"age"`
//		Friends []string `json:"friends"`
//	}
func GetTimeOut(ctx context.Context, requestUrl string, params interface{}, headers map[string]string, timeOut time.Duration, result interface{}) error {
	ctx = utils.WithChildSpan(ctx, "get:"+requestUrl)

	//解析验证url
	_, err := url.Parse(requestUrl)
	if err != nil {
		logger.SpanFailed(ctx, "parse url failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(err))
		return err
	}

	queryValues := url.Values{}
	if params != nil {
		queryValues, err = query.Values(params)
		if err != nil {
			logger.SpanFailed(ctx, "parse query params failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(err))
			return err
		}
	}

	logger.Info(ctx, "request url", zap.String("url", requestUrl), zap.Any("header", headers))

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	//请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	if params != nil {
		req.SetRequestURI(fmt.Sprintf("%s?%s", requestUrl, queryValues.Encode()))
	} else {
		req.SetRequestURI(fmt.Sprintf("%s", requestUrl))
	}

	req.Header.SetMethod(fasthttp.MethodGet)

	//请求数据
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = doTimeOut(ctx, req, resp, timeOut, result)
	if err != nil {
		logger.SpanFailed(ctx, "do request failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(err))
		return err
	}

	logger.SpanSuccess(ctx, "request success", zap.String("req", req.String()))
	return nil
}

// PostTimeOut Post请求，不需要设置超时时间
//
// params: 请求参数为结构体类型，需要在字段后面加 json tag，会自动转换为对应的参数
//
//	type request struct {
//		Name    string   `json:"name"`
//		Age     int      `json:"age"`
//		Friends []string `json:"friends"`
//	}
func PostTimeOut(ctx context.Context, requestUrl string, params interface{}, body interface{}, headers map[string]string, timeOut time.Duration, result interface{}) error {
	ctx = utils.WithChildSpan(ctx, "post:"+requestUrl)

	//解析验证url
	_, err := url.Parse(requestUrl)
	if err != nil {
		logger.SpanFailed(ctx, "parse url failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("body", body), zap.Any("header", headers), zap.Error(err))
		return err
	}

	//url 参数
	queryValues := url.Values{}
	if params != nil {
		queryValues, err = query.Values(params)
		if err != nil {
			logger.SpanFailed(ctx, "parse query params failed", zap.String("requestUrl", requestUrl), zap.Any("params", params), zap.Any("header", headers), zap.Error(err))
			return err
		}
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(requestUrl)

	//请求参数
	if params != nil {
		req.SetRequestURI(fmt.Sprintf("%s?%s", requestUrl, queryValues.Encode()))
	} else {
		req.SetRequestURI(fmt.Sprintf("%s", requestUrl))
	}

	//请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// header
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType(constant.ContentTypeJson)
	if headers[constant.HeaderKeyContextType] == constant.ContentTypeForm {
		req.Header.SetContentType(constant.ContentTypeForm)
		values, err := query.Values(body)
		if err != nil {
			logger.SpanFailed(ctx, "form data encode failed", zap.String("requestUrl", requestUrl), zap.Any("body", body), zap.Any("header", headers), zap.Error(err))
			return err
		}
		req.SetBodyString(values.Encode())
	} else {
		marshal, err := sonic.Marshal(body)
		if err != nil {
			logger.SpanFailed(ctx, "json marshal failed", zap.String("requestUrl", requestUrl), zap.Any("body", body), zap.Any("header", headers), zap.Error(err))
			return err
		}
		req.SetBody(marshal)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 请求并解析数据
	err = doTimeOut(ctx, req, resp, timeOut, result)
	if err != nil {
		return err
	}

	logger.SpanSuccess(ctx, "request success", zap.String("req", req.String()))
	return nil
}

func doTimeOut(ctx context.Context, req *fasthttp.Request, resp *fasthttp.Response, timeOut time.Duration, result interface{}) error {
	var err error

	if timeOut <= 0 {
		err = client.Do(req, resp)
	} else {
		err = client.DoTimeout(req, resp, timeOut)
	}

	if err != nil {
		logger.SpanFailed(ctx, "http get failed", zap.String("req", req.String()), zap.Error(err))
		return err
	}

	//状态码验证
	statusCode := resp.StatusCode()
	if statusCode != http.StatusOK {
		errInner := errors.New("data request failed , code:" + strconv.Itoa(statusCode))
		logger.SpanFailed(ctx, "http get failed", zap.Int("code", statusCode), zap.String("req", req.String()), zap.Error(errInner))
		return errInner
	}

	//反序列化参数
	err = sonic.Unmarshal(resp.Body(), result)
	if err != nil {
		logger.SpanFailed(ctx, "json unmarshal failed", zap.Int("code", statusCode), zap.String("req", req.String()), zap.String("data", string(resp.Body())), zap.Error(err))
		return err
	}

	return nil
}
