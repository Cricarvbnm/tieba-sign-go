package tieba

import (
	"sync"

	"tieba-sign/client"
	"tieba-sign/log"
)

type Signer struct {
	client    *client.Client
	failTotal int
	mu        sync.Mutex
}

func NewSigner(tbClient *client.Client) *Signer {
	return &Signer{client: tbClient}
}

func (s *Signer) ForumSign(forumName string) error {
	tbs, err := TBSFetch(s.client)
	if err != nil {
		s.mu.Lock()
		s.failTotal++
		s.mu.Unlock()
		return err
	}

	if err := ForumSign(s.client, forumName, tbs); err != nil {
		s.mu.Lock()
		s.failTotal++
		s.mu.Unlock()
		return err
	}

	log.Info.Println("签到完成:", forumName)
	return nil
}

func (s *Signer) FailTotal() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.failTotal
}
