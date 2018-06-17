package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

const (
	startUser  = "CroMarmot"
	workerSize = 100
	digDeep    = 2
	userfile   = "res/user.json"
	relfile    = "res/rel.json"
	mainlog    = "log/main.log"
)

/*
user -> user2url() -> url

worker(){
	rep(){
		pickOneUrl from fifo buffer and analyze{
		fifo buffer <- each
	}
}

}

*/

type CustomErr int

func (ce CustomErr) Error() string {
	return string(ce)
}

type Relation struct {
	Follower string `json:"fer"`
	Followee string `json:"fee"`
}

type StrDep struct {
	User string
	Url  string
	Dep  int
}

type CrawlCtrl interface {
	FiUrl(strDep StrDep)
	FoUrl() (StrDep, error)
	Worker()
	Init()
}

type CrawlManager struct {
	muxfri   sync.Mutex
	friends  map[string]bool
	rels     []Relation
	muxurl   sync.Mutex
	urls     map[string]bool
	urlsfifo *list.List
	logger   *log.Logger
	wg       sync.WaitGroup
}

func NewCrawlManager() *CrawlManager {
	um := new(CrawlManager)
	f, err := os.OpenFile(mainlog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	um.logger = log.New(f, "MAIN:\t", log.Lshortfile)
	um.friends = make(map[string]bool)
	um.urls = make(map[string]bool)
	um.rels = make([]Relation, 0)
	um.urlsfifo = list.New()
	return um
}

func (um *CrawlManager) Init() {
	um.wg.Add(workerSize)
	um.FiUrl(StrDep{startUser, user2url(startUser), digDeep})
	for i := 0; i < workerSize; i++ {
		go um.Worker()
	}
}

func user2url(uname string) string {
	return "https://github.com/" + uname + "?tab=followers"
}

func (um *CrawlManager) AddFriend(followee string, follower string) bool {
	um.muxfri.Lock()
	defer um.muxfri.Unlock()
	_, ok := um.friends[follower]

	um.rels = append(um.rels, Relation{follower, followee})

	if !ok {
		um.friends[follower] = true
	}
	return !ok
}

func (um *CrawlManager) FiUrl(strDep StrDep) {
	um.wg.Add(1)
	um.muxurl.Lock()
	if _, ok := um.urls[strDep.Url]; !ok {
		um.urls[strDep.Url] = true
		um.urlsfifo.PushBack(strDep)
	} else {
		um.wg.Done()
	}
	um.muxurl.Unlock()
}

func (um *CrawlManager) getBody(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		um.logger.Print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		um.logger.Print(err)
	}
	return string(body)
}

func (um *CrawlManager) FoUrl() (StrDep, error) {
	waitCnt := 0
	for {
		um.muxurl.Lock()
		if v := um.urlsfifo.Front(); v != nil {
			defer um.muxurl.Unlock()
			um.urlsfifo.Remove(v)
			return v.Value.(StrDep), nil
		}
		um.muxurl.Unlock()
		waitCnt++
		time.Sleep(100 * time.Millisecond)
		if waitCnt > 2 {
			return StrDep{}, CustomErr(-1)
		}
	}
}

func (um *CrawlManager) analyzeUrl(strDep StrDep) {
	if strDep.Dep == 0 {
		return
	}
	s := um.getBody(strDep.Url)
	//fmt.Print(s)

	// friends
	pattern1 := `<span class="link-gray pl-1">(.*?)</span>`
	rp1 := regexp.MustCompile(pattern1)
	finddata1 := rp1.FindAllStringSubmatch(s, -1)
	//fmt.Printf("--------------------Friends\n")
	//fmt.Print(finddata1)
	//fmt.Printf("--------------------Friends\n")
	for _, v := range finddata1 {
		newUser := v[1]
		um.AddFriend(strDep.User, newUser)
		um.logger.Printf("REL: %v - %v\n", strDep.User, newUser)
		um.FiUrl(StrDep{newUser, user2url(newUser), strDep.Dep - 1})
	}

	// next page
	patternnext := `<a rel="nofollow" href="(.*?)">Next</a>`
	rp2 := regexp.MustCompile(patternnext)
	finddata2 := rp2.FindAllStringSubmatch(s, -1)
	if len(finddata2) == 1 {
		um.FiUrl(StrDep{strDep.User, finddata2[0][1], strDep.Dep - 1})
	}
}

func (um *CrawlManager) Worker() {
	errcnt := 0
	for {
		if errcnt > 3 {
			break
		}
		v, err := um.FoUrl()
		if err != nil {
			errcnt++
			time.Sleep(100 * time.Millisecond)
			continue
		}
		errcnt = 0
		um.analyzeUrl(v)
		um.wg.Done()
	}
	um.wg.Done()
}

func main() {
	var crawlCtrl CrawlCtrl
	crawlCtrl = NewCrawlManager()
	crawlCtrl.Init()
	crawlCtrl.(*CrawlManager).wg.Wait()
	crawlCtrl.(*CrawlManager).logger.Println("Fetched all")

	js := NewJsonSaver()
	js.Save(userfile, crawlCtrl.(*CrawlManager).friends)
	crawlCtrl.(*CrawlManager).logger.Println("um.rels{{{")
	crawlCtrl.(*CrawlManager).logger.Println(js.SaveMem(crawlCtrl.(*CrawlManager).rels))
	crawlCtrl.(*CrawlManager).logger.Println("}}}um.rels")
	js.Save(relfile, crawlCtrl.(*CrawlManager).rels)
	// crawlCtrl.(*CrawlManager).logger.Println("Saved all")
}
