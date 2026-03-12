package tieba

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"tieba-sign/client"
	"tieba-sign/log"
)

const (
	forumsURL = "https://tieba.baidu.com/mo/q/newmoindex"
	tbsURL    = "https://tieba.baidu.com/dc/common/tbs"
)

func ForumsFetch(tbClient *client.Client) ([]Forum, error) {
	respBody, err := tbClient.Get(forumsURL)
	if err != nil {
		return nil, fmt.Errorf("无法获取关注贴吧列表: %w", err)
	}

	if err := log.ToFile("forums.json", respBody); err != nil {
		return nil, err
	}

	resp := struct {
		Error string `json:"error"`
		Data  struct {
			Forums []Forum `json:"like_forum"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("无法解析关注贴吧列表: %w", err)
	}

	if resp.Error != "success" {
		return nil, fmt.Errorf("无法获取关注贴吧列表: %s", resp.Error)
	}

	return resp.Data.Forums, nil
}

func TBSFetch(tbClient *client.Client) (string, error) {
	respBody, err := tbClient.Get(tbsURL)
	if err != nil {
		return "", fmt.Errorf("无法获取 TBS: %w", err)
	}

	if err := log.ToFile("tbs.json", respBody); err != nil {
		return "", err
	}

	resp := struct {
		TBS     string `json:"tbs"`
		IsLogin int    `json:"is_login"`
	}{}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("无法解析 TBS: %w", err)
	}

	if resp.IsLogin == 0 {
		return "", fmt.Errorf("未登录")
	}

	return resp.TBS, nil
}

func ForumSign(tbClient *client.Client, forumName string, tbs string) error {
	const signURL = "https://tieba.baidu.com/sign/add"

	reqBody := url.Values{
		"kw":  {forumName},
		"tbs": {tbs},
		"ie":  {"utf-8"},
	}.Encode()

	logDir := "sign-forum/" + forumName
	if err := log.ToFile(logDir+"/req-body", []byte(reqBody)); err != nil {
		return err
	}

	respBody, err := tbClient.Post(signURL, strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("签到失败: %s: %w", forumName, err)
	}

	if err := log.ToFile(logDir+"/resp-body.json", respBody); err != nil {
		return err
	}

	resp := struct {
		Error string `json:"error"`
	}{}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("签到失败: %s: %w", forumName, err)
	}

	if resp.Error != "" {
		return fmt.Errorf("签到失败: %s: %s", forumName, resp.Error)
	}

	return nil
}
