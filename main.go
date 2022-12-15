package main

import (
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/CesarDelgadoM/webdozer/crawler"
	"github.com/CesarDelgadoM/webdozer/database"
	"github.com/CesarDelgadoM/webdozer/store"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatal("No arguments...")
	}

	url := os.Args[1]
	patterns := os.Args[2]

	procs := runtime.GOMAXPROCS(runtime.NumCPU())
	redis := database.NewRedisPool("localhost:6379", procs)
	defer redis.Close()

	var wg sync.WaitGroup

	wg.Add(2)
	c := crawler.NewCrawler(&wg, redis)
	c.LaunchCrawler(url)

	s := store.NewStore(url, &wg, redis)
	s.StoreUrls(c.UrlsFound, patterns)
	wg.Wait()

	log.Println("Finished the process")
}
