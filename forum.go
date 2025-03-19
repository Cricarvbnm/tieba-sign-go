package main

import (
	"fmt"
	"regexp"
)

func getTbs(forumName string) (string, error) {
	html, err := fetchHtml(fmt.Sprintf("https://tieba.baidu.com/f?kw=%s", forumName))
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`['"]tbs['"]\s*:\s*['"]([^'"]+)['"]`)
	matches := re.FindStringSubmatch(string(html))
	if matches == nil {
		forumsLogger.Printf("<!-- %s -->\n%s\n", forumName, html)
		return "", fmt.Errorf("解析tbs失败")
	}

	return matches[1], nil
}
