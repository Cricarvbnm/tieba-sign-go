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
			log.Notice.Println("没有关注的吧")
		} else {
			log.Notice.Println("没有需要签到的吧")
		}
		return
	}

	log.Notice.Printf("待签到: %d/%d\n", unsigned, total)

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

	log.Notice.Printf("签到完成, 失败: %d/%d\n", signer.FailTotal(), unsigned)
}

func configLoad() *config.Config {
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		log.Fatal.Fatalln(err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal.Fatalln(err)
	}
	return cfg
}

func forumsFetch(tbClient *client.Client) []tieba.Forum {
	forums, err := tieba.ForumsFetch(tbClient)
	if err != nil {
		log.Fatal.Fatalln(err)
	}
	return forums
}

func forumSign(signer *tieba.Signer, forumName string) {
	if err := signer.ForumSign(forumName); err != nil {
		log.Err.Println(err)
		return
	}
}
