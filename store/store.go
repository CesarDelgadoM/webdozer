package store

import (
	"log"
	"sync"

	"github.com/CesarDelgadoM/webdozer/database"
	"github.com/CesarDelgadoM/webdozer/utils"
)

const (
	SEARCH_WORDS string = "search-words"
)

type Store struct {
	url  string
	list database.Repository
	hash database.Hash
	wg   *sync.WaitGroup
}

func NewStore(url string, wg *sync.WaitGroup, redis *database.RedisPool) *Store {

	return &Store{
		url:  url,
		wg:   wg,
		list: database.NewList(redis),
		hash: database.NewHash(SEARCH_WORDS, redis),
	}
}

func (s *Store) StoreUrls(urls chan string, patterns string) {

	storeKey := "urls-" + utils.ExtractNameUrl(s.url)

	s.hash.Set(storeKey, patterns)

	go func() {
		defer s.wg.Done()

		for url := range urls {
			log.Println(url)
			err := s.list.Add(storeKey, url)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}
