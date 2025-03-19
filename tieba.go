package main

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func parseForums(html string) ([]Forum, error) {
	re := regexp.MustCompile(`['"]forums['"]\s*:\s*(\[[^\]]*\])`)
	matches := re.FindStringSubmatch(html)
	if matches == nil {
		homepageLogger.Println(html)
		return nil, fmt.Errorf("解析贴吧列表失败")
	}

	var forums []Forum
	err := json.Unmarshal([]byte(matches[1]), &forums)
	if err != nil {
		homepageLogger.Println(html)
		return nil, fmt.Errorf("解析贴吧列表失败: %w", err)
	}

	return forums, nil
}
