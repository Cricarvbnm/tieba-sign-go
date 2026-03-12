package tieba

import (
	"sync"

	"tieba-sign/client"
	"tieba-sign/log"
)

type Signer struct {
	client       *client.Client
	succeedTotal int
	mu           sync.Mutex
}

func NewSigner(tbClient *client.Client) *Signer {
	return &Signer{client: tbClient}
}

func (s *Signer) ForumSign(forumName string) error {
	tbs, err := TBSFetch(s.client)
	if err != nil {
		return err
	}

	if err := ForumSign(s.client, forumName, tbs); err != nil {
		return err
	}

	s.mu.Lock()
	s.succeedTotal++
	s.mu.Unlock()

	log.Info.Println("签到完成:", forumName)
	return nil
}

func (s *Signer) SucceedTotal() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.succeedTotal
}
