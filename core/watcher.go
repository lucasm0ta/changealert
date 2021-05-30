package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"sort"
	"sync"
	"time"
)

type Watched struct {
	// Local identification of the client
	URL        *url.URL  `json:"url"`
	Hash       string    `json:"hash"`
	LastUpdate time.Time `json:"last_update"`
	Interval   int64     `json:"interval"`
}

type WatcherService struct {
	WatchedUrls []*Watched
	httpService *HttpService
	repository  *WatcherRepository
}

func NewWatcherService() *WatcherService {
	watcher := new(WatcherService)
	watcher.WatchedUrls = nil
	watcher.httpService = NewHttpService()
	watcher.repository = NewWatcherRepository()
	return watcher
}

func (watcherService *WatcherService) Start() {
	for {
		var wg sync.WaitGroup
		for _, watched := range watcherService.WatchedUrls {
			current := time.Now().Unix()
			nextUpdate := watched.LastUpdate.Unix() + watched.Interval
			if nextUpdate < current {
				wg.Add(1)
				go watcherService.Request(&wg, watched)
			}
		}
		wg.Wait()
	}
}

func (watcherService *WatcherService) Request(wg *sync.WaitGroup, watched *Watched) {
	defer wg.Done()
	body, _ := watcherService.httpService.GetBody(watched.URL.String())
	hash := md5.New()
	io.WriteString(hash, body)
	watched.Hash = string(hash.Sum(nil))

	fmt.Println(watched.URL)
	watched.LastUpdate = time.Now()
}

func (watcherService *WatcherService) Register(request *WatchRequest) error {
	watched := new(Watched)
	watched.URL = request.URL
	watched.Interval = 10

	body, err := watcherService.httpService.GetBody(watched.URL.String())
	if err != nil {
		return err
	}
	hash := md5.New()
	io.WriteString(hash, body)

	watched.Hash = string(hash.Sum(nil))
	watched.LastUpdate = time.Now()

	fmt.Println(watched.URL)
	watcherService.Insert(watched)
	return nil
}

func (watcherService *WatcherService) Insert(watch *Watched) {
	watcherService.persistWatched(watch)
	isAfterIndex := func(i int) bool { return watcherService.WatchedUrls[i].LastUpdate.After(watch.LastUpdate) }
	idx := sort.Search(len(watcherService.WatchedUrls), isAfterIndex)
	watcherService.WatchedUrls = append(watcherService.WatchedUrls, nil)
	copy(watcherService.WatchedUrls[idx+1:], watcherService.WatchedUrls[idx:])
	watcherService.WatchedUrls[idx] = watch
}

func (watcherService *WatcherService) persistWatched(watch *Watched) {
	watcherService.repository
}
