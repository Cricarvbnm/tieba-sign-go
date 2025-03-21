package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	client *Client

	domainURL = "https://tieba.baidu.com"
)

type Client struct {
	client *http.Client
	header http.Header
}

func (c *Client) Get(url string) ([]byte, error) {
	return c.fetch("GET", url, nil)
}

func (c *Client) Post(url string, body io.Reader) ([]byte, error) {
	return c.fetch("POST", url, body)
}

func (c *Client) fetch(method string, url string, body io.Reader) ([]byte, error) {
	// 创建请求
	req, err := newTiebaRequest(method, url, body, c.header)
	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %w", err)
	}

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %s", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败: %s", resp.Status)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取响应体: %w", err)
	}
	return bodyBytes, nil
}

func initClient() {
	// cookie
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	domainURLStruct, _ := url.Parse(domainURL)
	cookieJar.SetCookies(domainURLStruct, []*http.Cookie{
		{
			Name:  "BDUSS",
			Value: config.BDUSS,
			Path:  "/",
		},
		{
			Name:  "STOKEN",
			Value: config.STOKEN,
			Path:  "/",
		},
	})

	// client
	httpClient := &http.Client{Jar: cookieJar}

	// header
	header := http.Header{
		"Host":       {"tieba.baidu.com"},
		"User-Agent": {`Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0`},
	}

	// construct
	client = &Client{
		client: httpClient,
		header: header,
	}
}

func newTiebaRequest(method string, url string, body io.Reader, additionalHeader http.Header) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	for k, v := range additionalHeader {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	return req, err
}
