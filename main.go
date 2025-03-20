package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	homepageLogger *log.Logger
	forumsLogger   *log.Logger
)

var (
	failedCount = 0
)

func main() {
	// get homepage
	homepageHTML, err := fetchHomepageHTML()
	if err != nil {
		log.Fatalln(err)
	}

	// get forums
	forums, err := parseForums(homepageHTML)
	if err != nil {
		log.Fatalln(err)
	} else if len(forums) == 0 {
		log.Println("warning: 没有找到关注的吧")
	}

	// get unsigned forums
	unsignedForums := filterUnsignedForums(forums)

	// statistic
	forumCount := len(forums)
	unsignedForumsCount := len(unsignedForums)

	if unsignedForumsCount == 0 {
		log.Println("没有需要签到的吧")
		return
	} else {
		log.Printf("待签到: %d/%d\n", unsignedForumsCount, forumCount)
	}

	// sign
	for _, forum := range unsignedForums {
		sign(forum.Name)
		time.Sleep(1 * time.Second)
	}

	log.Printf("签到完成(失败/总签到数): %d/%d\n", failedCount, unsignedForumsCount)
}

func init() {
	initConfig()
	initLogger()
	initRequest()
}

func initLogger() {
	log.SetFlags(0)

	homepageLogger = log.New(openLogFile("homepage.html"), "", 0)
	forumsLogger = log.New(openLogFile("forums.json"), "", 0)
}

func openLogFile(path string) *os.File {
	// log path
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}

	logDir := filepath.Join(dataHome, "tieba-sign", "log")
	path = filepath.Join(logDir, path)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatalf("Failed to create log directory: %s: %s\n", logDir, err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Failed to open %s: %s\n", path, err)
	}
	return file
}

func sign(forumName string) {
	log.Println("正在签到:", forumName)

	tbs, err := fetchForumTbs(forumName)
	if err != nil {
		failedCount++
		log.Println("获取 tbs 失败:", err)
		return
	}

	if _, err := signForum(forumName, tbs, false); err != nil {
		failedCount++
		log.Println("签到失败:", err)
	} else {
		log.Println("签到完成:", forumName)
	}
}
