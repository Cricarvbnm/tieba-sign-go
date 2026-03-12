package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	domainURL = "https://tieba.baidu.com"
)

type Client struct {
	httpClient *http.Client
	header     http.Header
}

func New(bduss, stoken string) *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	domainURLStruct, _ := url.Parse(domainURL)
	cookieJar.SetCookies(domainURLStruct, []*http.Cookie{
		{Name: "BDUSS", Value: bduss, Path: "/"},
		{Name: "STOKEN", Value: stoken, Path: "/"},
	})

	httpClient := &http.Client{Jar: cookieJar}

	header := http.Header{
		"Host":       {"tieba.baidu.com"},
		"User-Agent": {`Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0`},
	}

	return &Client{
		httpClient: httpClient,
		header:     header,
	}
}

func (c *Client) Get(url string) ([]byte, error) {
	return c.fetch("GET", url, nil)
}

func (c *Client) Post(url string, body io.Reader) ([]byte, error) {
	return c.fetch("POST", url, body)
}

func (c *Client) fetch(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %w", err)
	}

	for k, v := range c.header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取响应体: %w", err)
	}

	return bodyBytes, nil
}
