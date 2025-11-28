package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type HTTPClient struct {
	client *http.Client
	domain string
}

func NewHTTPClient(domain string) (*HTTPClient, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain cannot be empty")
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: 15 * time.Second, // 你可以调整
		},
		domain: domain,
	}, nil
}

func (d *HTTPClient) commonRequest(
	method string,
	requestPath string,
	query map[string]string,
	header map[string]string,
	body any,
) (*http.Response, error) {

	// 拼 URL
	u, err := url.Parse(d.domain)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, requestPath)

	// query 参数
	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	// 序列化 body
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body failed: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	// 创建 request
	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置 header
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 发起请求
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
