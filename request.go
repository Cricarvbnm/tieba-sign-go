package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	client *http.Client

	domain = "https://tieba.baidu.com"
)

func initRequest() {
	cookies, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client = &http.Client{Jar: cookies}

	domainURL, _ := url.Parse(domain)
	if config.Cookie == "" {
		cookies.SetCookies(domainURL, []*http.Cookie{
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
	}
}

func makeTiebaHeader() http.Header {
	header := http.Header{
		"Host":       []string{"tieba.baidu.com"},
		"User-Agent": []string{config.UserAgent},
	}
	if config.Cookie != "" {
		header.Add("Cookie", config.Cookie)
	}
	return header
}

func makeTiebaRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}

	req.Header = makeTiebaHeader()
	return req
}

func fetch(method string, url string, body io.Reader) ([]byte, error) {
	req := makeTiebaRequest(method, url, body)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status: %s", resp.Status)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, nil
}
