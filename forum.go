package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
)

var (
	signURL       = domain + "/sign/add"
	signURLMobile = "https://c.tieba.baidu.com/c/c/forum/sign"
	// forumURL      = domain + "/f"
)

func filterUnsignedForums(forums []Forum) []Forum {
	var filteredForums []Forum
	for _, forum := range forums {
		if forum.IsSign == 0 {
			filteredForums = append(filteredForums, forum)
		}
	}
	return filteredForums
}

func fetchForumTbs(forumName string) (string, error) {
	tbsURL := "https://tieba.baidu.com/dc/common/tbs"
	tbsJSONBytes, err := fetch("GET", tbsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch tbs: %w", err)
	}

	tbsStruct := struct {
		Tbs     string `json:"tbs"`
		IsLogin int    `json:"is_login"`
	}{}
	if err := json.Unmarshal(tbsJSONBytes, &tbsStruct); err != nil {
		return "", fmt.Errorf("failed to unmarshal tbs: %w", err)
	} else if tbsStruct.IsLogin == 0 {
		return "", fmt.Errorf("not logged in")
	}

	return tbsStruct.Tbs, nil
}

func parseTbs(forumHTML string) (string, error) {
	re := regexp.MustCompile(`['"]tbs['"]\s*:\s*['"](.*?)['"]`)
	matches := re.FindStringSubmatch(forumHTML)
	if matches == nil {
		return "", fmt.Errorf("failed to parse tbs")
	}

	return matches[1], nil
}

func logForumHTML(forumName string, forumHTML string) {
	file := openLogFile(fmt.Sprintf("forum-html/%s.html", forumName))
	defer file.Close()
	forumsLogger := log.New(file, "", 0)
	forumsLogger.Println(forumHTML)
}

func signForum(forumName string, tbs string, isSignMobile bool) (string, error) {
	// make request
	var (
		bodyBytes []byte
		err       error
	)

	if isSignMobile {
		bodyBytes, err = signMobile(forumName, tbs)
	} else {
		bodyBytes, err = signPC(forumName, tbs)
	}

	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	var bodyJSON struct {
		ErrorMsg string `json:"error_msg"`
	}
	json.Unmarshal(bodyBytes, &bodyJSON)
	if bodyJSON.ErrorMsg != "" {
		return "", fmt.Errorf("failed to sign: %s", bodyJSON.ErrorMsg)
	}

	// read
	body := string(bodyBytes)

	// log
	logSignResp(forumName, body)

	return body, nil
}

func signPC(forumName string, tbs string) ([]byte, error) {
	data := url.Values{
		"kw":  {forumName},
		"tbs": {tbs},
		"ie":  {"utf-8"},
	}.Encode()

	logSignPost(forumName, string(data))

	return fetch("POST", signURL, strings.NewReader(data))
}

func signMobile(forumName string, tbs string) ([]byte, error) {
	signStr := md5.Sum(fmt.Appendf(nil, `kw=%s&tbs=%stiebaclient!!!`, forumName, tbs))

	data := url.Values{
		"kw":   {forumName},
		"tbs":  {tbs},
		"sign": {fmt.Sprintf("%x", signStr)},
	}.Encode()

	logSignPost(forumName, data)

	return fetch("POST", signURLMobile, strings.NewReader(data))
}

func logSignResp(forumName string, respBody string) {
	file := openLogFile(fmt.Sprintf("sign-resp/%s.json", forumName))
	defer file.Close()
	signRespLogger := log.New(file, "", 0)
	signRespLogger.Println(respBody)
}

func logSignPost(forumName string, postData string) {
	file := openLogFile(fmt.Sprintf("sign-post/%s.json", forumName))
	defer file.Close()
	signPostLogger := log.New(file, "", 0)
	signPostLogger.Println(postData)
}
