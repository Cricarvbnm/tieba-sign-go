package main

import (
	"sync"

	"tieba-sign/client"
	"tieba-sign/config"
	"tieba-sign/log"
	"tieba-sign/tieba"
)

func main() {
	cfg := configLoad()
	tbClient := client.New(cfg.BDUSS, cfg.STOKEN)

	forums := forumsFetch(tbClient)
	total := len(forums)

	forums = tieba.ForumsFilterUnsigned(forums)
	unsigned := len(forums)

	if unsigned == 0 {
		if total == 0 {
			log.Warn.Println("没有关注的吧")
		} else {
			log.Warn.Println("没有需要签到的吧")
		}
		return
	}

	log.Info.Printf("待签到: %d/%d\n", unsigned, total)

	signer := tieba.NewSigner(tbClient)
	taskChan := make(chan string, unsigned)

	var wg sync.WaitGroup
	wg.Add(unsigned)

	go func() {
		for _, forum := range forums {
			taskChan <- forum.Name
		}
		close(taskChan)
	}()

	for i := 0; i < 4; i++ {
		go func() {
			for forumName := range taskChan {
				forumSign(signer, forumName)
				wg.Done()
			}
		}()
	}

	wg.Wait()

	log.Info.Printf("签到完成, 成功: %d/%d\n", signer.SucceedTotal(), unsigned)
}

func configLoad() *config.Config {
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		log.Error.Fatalln(err)
	}
	if err := cfg.Validate(); err != nil {
		log.Error.Fatalln(err)
	}
	return cfg
}

func forumsFetch(tbClient *client.Client) []tieba.Forum {
	forums, err := tieba.ForumsFetch(tbClient)
	if err != nil {
		log.Error.Fatalln(err)
	}
	return forums
}

func forumSign(signer *tieba.Signer, forumName string) {
	log.Info.Println("正在签到:", forumName)

	if err := signer.ForumSign(forumName); err != nil {
		log.Error.Println(err)
		return
	}
}
