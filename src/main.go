package main

import (
	"log"
	"sync"
)

var (
	forumCount         int
	unsignedForumCount int
	succeedCount       = 0

	signWaitGroup = sync.WaitGroup{}
	signTasks     = make(chan string)
)

func main() {
	// get forums
	forums, err := getForums()
	if err != nil {
		ErrorLogger.Fatalln(err)
	}
	forumCount = len(forums)
	if forumCount == 0 {
		WarnLogger.Println("没有关注的吧")
		return
	}

	// filter forums
	forums = filterUnsignedForums(forums)
	unsignedForumCount = len(forums)
	if unsignedForumCount == 0 {
		WarnLogger.Println("没有需要签到的吧")
		return
	}

	// statistic
	log.Printf("未签/总关注数: %d/%d\n", unsignedForumCount, forumCount)

	// task provider goroutine
	signWaitGroup.Add(unsignedForumCount)
	go func() {
		for _, forum := range forums {
			signTasks <- forum.Name
		}
	}()

	// task worker goroutines
	for range 4 {
		go func() {
			for forumName := range signTasks {
				sign(forumName)
			}
		}()
	}

	signWaitGroup.Wait()

	log.Printf("签到完成, 成功/总签到数: %d/%d\n", succeedCount, unsignedForumCount)
}

func init() {
	if err := loadConfig(); err != nil {
		ErrorLogger.Fatalln(err)
	}
	initLog()
	initClient()
}

func filterUnsignedForums(forums []Forum) []Forum {
	var unsignedForums []Forum
	for _, forum := range forums {
		if forum.IsSign == 0 {
			unsignedForums = append(unsignedForums, forum)
		}
	}
	return unsignedForums
}

func sign(forumName string) {
	defer signWaitGroup.Done()

	log.Println("正在签到:", forumName)

	// get tbs
	tbs, err := getTBS()
	if err != nil {
		ErrorLogger.Println(err)
		return
	}

	// sign
	if err := signForum(forumName, tbs); err != nil {
		ErrorLogger.Println(err)
		return
	}

	succeedCount++
	log.Println("签到完成:", forumName)
}
