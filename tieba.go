package main

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Forum struct {
	Name   string `json:"forum_name"`
	IsSign int    `json:"is_sign"`
}

type TBS struct {
	Tbs     string `json:"tbs"`
	IsLogin int    `json:"is_login"`
}

func fetchHomepageHTML() (string, error) {
	// fetch
	homepageHTMLBytes, err := fetch("GET", domain, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch homepage: %w", err)
	}

	// read
	homepageHTML := string(homepageHTMLBytes)
	homepageLogger.Println(homepageHTML)

	// check
	if isSecurityCheck(homepageHTML) {
		return "", fmt.Errorf("failed to fetch homepage: 百度安全验证")
	}

	return homepageHTML, nil
}

func isSecurityCheck(html string) bool {
	re := regexp.MustCompile(`<title>百度安全验证</title>`)
	matches := re.FindStringIndex(html)
	return matches != nil
}

func parseForums(html string) ([]Forum, error) {
	// get json string
	re := regexp.MustCompile(`['"]forums['"']\s*:\s*(\[.*?\])`)
	matches := re.FindStringSubmatch(html)
	if matches == nil {
		return nil, fmt.Errorf("failed to parse forums")
	}
	defer forumsLogger.Println(matches[1])

	// parse json
	var forums []Forum
	err := json.Unmarshal([]byte(matches[1]), &forums)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal forums: %w", err)
	}

	return forums, nil
}
