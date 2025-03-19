package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	config       *Config
	client       *http.Client
	waitGroup    *sync.WaitGroup
	semaphore    chan struct{}
	successCount int = 0
)

var (
	homepageLogger *log.Logger
	forumsLogger   *log.Logger
	signPostLogger *log.Logger
)

func init() {
	// log
	log.SetFlags(0)

	homepageFile, err := os.OpenFile("log/homepage.log.html", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	homepageLogger = log.New(homepageFile, "", 0)

	forumsFile, err := os.OpenFile("log/forums.log.html", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	forumsLogger = log.New(forumsFile, "", 0)

	signPostFile, err := os.OpenFile("log/sign-post.log.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	signPostLogger = log.New(signPostFile, "", 0)

	// config
	config, err = loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// request
	waitGroup = &sync.WaitGroup{}
	semaphore = make(chan struct{}, config.Request.Concurrency)

	//client
	client = &http.Client{}
}

func main() {
	html, err := fetchHtml("https://tieba.baidu.com")
	if err != nil {
		log.Fatalln(err)
	} else if isSecurityVerification(html) {
		log.Fatalln("Cookie失效, 请重新获取, 已触发百度安全验证")
	}

	forums, err := parseForums(html)
	if err != nil {
		log.Fatalln(err)
	} else if len(forums) == 0 {
		log.Fatalln("贴吧列表为空")
	}

	var unsignedForums []Forum
	for _, forum := range forums {
		if forum.IsSign == 0 {
			unsignedForums = append(unsignedForums, forum)
		}
	}

	log.Printf("还需签到贴吧数量: %d\n", len(unsignedForums))
	for _, forum := range unsignedForums {
		waitGroup.Add(1)
		semaphore <- struct{}{}
		sign(forum.Name) // 暂时不用并发

		time.Sleep(1 * time.Second) // 避免请求过快
	}
	waitGroup.Wait()

	log.Printf("签到成功统计: %d/%d\n", successCount, len(unsignedForums))
}

type Forum struct {
	Name   string `json:"forum_name"`
	IsSign int    `json:"is_sign"`
}

func isSecurityVerification(html string) bool {
	re := regexp.MustCompile(`<title>百度安全验证</title>`)
	matches := re.FindStringSubmatch(html)
	return len(matches) > 0
}

func sign(forumName string) {
	defer waitGroup.Done()
	defer func() { <-semaphore }()
	defer signPostLogger.Println()

	tbs, err := getTbs(forumName)
	if err != nil {
		log.Println(fmt.Errorf("获取tbs失败: %w", err))
		doSignFailed(forumName)
	}

	req := makeTiebaRequest(
		"POST",
		"https://tieba.baidu.com/sign/add",
		strings.NewReader(fmt.Sprintf(`{"kw": %s, "ie": %s, "tbs": %s}`, forumName, "utf-8", tbs)),
	)
	signPostLogger.Println(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(fmt.Errorf("请求失败: %w", err))
		doSignFailed(forumName)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("状态码异常: %d", resp.StatusCode)
		doSignFailed(forumName)
	}

	successCount++
	log.Printf("签到成功: %s\n", forumName)
}

func doSignFailed(forumName string) {
	log.Printf("签到失败: %s\n", forumName)
}
