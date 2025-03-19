package main

import (
	"fmt"
	"io"
	"net/http"
)

func makeTiebaRequest(method string, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Cookie", config.Cookie)
	req.Header.Set("Host", "tieba.baidu.com")
	return req
}

func fetchHtml(url string) (string, error) {
	req := makeTiebaRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("状态码异常: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	return string(buf), nil
}
