package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

var (
	forumsURL = "https://tieba.baidu.com/mo/q/newmoindex"
	tbsURL    = "https://tieba.baidu.com/dc/common/tbs"
	signURL   = "https://tieba.baidu.com/sign/add"
)

type Forum struct {
	Name   string `json:"forum_name"`
	IsSign int    `json:"is_sign"`
}

func getForums() ([]Forum, error) {
	// get response
	forumsRespBody, err := client.Get(forumsURL)
	if err != nil {
		return nil, err
	}
	if err := logToFile("forums.json", forumsRespBody); err != nil {
		return nil, err
	}

	// parse
	forumsResp := struct {
		Error string `json:"error"`
		Data  struct {
			Forums []Forum `json:"like_forum"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(forumsRespBody, &forumsResp); err != nil {
		return nil, fmt.Errorf("无法解析关注贴吧列表: %w", err)
	} else if forumsResp.Error != "success" {
		return nil, fmt.Errorf("无法获取关注贴吧列表: %s", forumsResp.Error)
	}

	return forumsResp.Data.Forums, nil
}

func getTBS() (string, error) {
	// get response
	tbsRespBody, err := client.Get(tbsURL)
	if err != nil {
		return "", fmt.Errorf("无法获取TBS: %w", err)
	}
	if err := logToFile("tbs.json", tbsRespBody); err != nil {
		return "", err
	}

	// parse
	tbsResp := struct {
		TBS     string `json:"tbs"`
		IsLogin int    `json:"is_login"`
	}{}
	if err := json.Unmarshal(tbsRespBody, &tbsResp); err != nil {
		return "", fmt.Errorf("无法解析TBS: %w", err)
	} else if tbsResp.IsLogin == 0 {
		return "", fmt.Errorf("未登录")
	}

	return tbsResp.TBS, nil
}

func signForum(forumName string, tbs string) error {
	// make reqBody
	reqBody := url.Values{
		"kw":  {forumName},
		"tbs": {tbs},
		"ie":  {"utf-8"},
	}.Encode()

	logDir := "sign-forum/" + forumName
	if err := logToFile(logDir+"/req-body", []byte(reqBody)); err != nil {
		return err
	}

	// request
	respBody, err := client.Post(signURL, strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("签到失败: %w", err)
	}
	if err := logToFile(logDir+"/resp-body.json", respBody); err != nil {
		return err
	}

	// check
	resp := struct {
		Error string `json:"error"`
	}{}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("签到失败: %w", err)
	} else if resp.Error != "" {
		return fmt.Errorf("签到失败: %s", resp.Error)
	}

	return nil
}
