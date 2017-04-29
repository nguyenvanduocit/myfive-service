package main

import (
	"fmt"
	"github.com/NYTimes/gziphandler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/github"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/medium"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/oxfordlearnersdictionaries"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/producthunt"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/rss"
	"github.com/nguyenvanduocit/myfive-crawler/interface"
	"github.com/nguyenvanduocit/myfive-service/config"
	"time"
)

type Post struct {
	ID    string `jsonapi:"primary,news"`
	Title string `jsonapi:"attr,title"`
	Url   string `jsonapi:"attr,url"`
}

type Site struct {
	ID      int    `jsonapi:"primary,sites"`
	Title   string `jsonapi:"attr,title"`
	Icon    string `jsonapi:"attr,icon"`
	Url     string `jsonapi:"attr,url"`
	FeedUrl string
	Crawler string
	Posts   []*Post `jsonapi:"relation,posts"`
}

type Server struct {
	Config *config.Config
	Sites  []*Site
}

func NewServer(config *config.Config) *Server {
	return &Server{
		Config: config,
		Sites: []*Site{
			{
				Title:   "A List Apart",
				Icon:    "alistapart.com.svg",
				Url:     "http://alistapart.com",
				ID:      0,
				FeedUrl: "http://alistapart.com/main/feed",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Toptal",
				Icon:    "toptal.com.png",
				Url:     "https://www.toptal.com",
				ID:      1,
				FeedUrl: "https://www.toptal.com/blog.rss",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Smashing Magazine",
				Url:     "https://www.smashingmagazine.com",
				ID:      2,
				FeedUrl: "https://www.smashingmagazine.com/feed/",
				Icon:    "smashingmagazine.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "David Walsh Blog",
				Url:     "https://davidwalsh.name",
				ID:      3,
				FeedUrl: "https://davidwalsh.name/feed",
				Icon:    "davidwalsh.name.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "SitePoint",
				Url:     "https://www.sitepoint.com",
				ID:      4,
				FeedUrl: "https://www.sitepoint.com/feed",
				Icon:    "sitepoint.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "GitHub Trending",
				ID:      5,
				Url:     "https://github.com/trending",
				FeedUrl: "https://github.com/trending",
				Icon:    "github.svg",
				Crawler: "github",
				Posts:   []*Post{},
			},
			{
				Title:   "Echo JS",
				ID:      6,
				Url:     "http://www.echojs.com",
				FeedUrl: "http://www.echojs.com/rss",
				Icon:    "echojs.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Coding Horror",
				ID:      7,
				Url:     "https://blog.codinghorror.com",
				FeedUrl: "https://blog.codinghorror.com/rss/",
				Icon:    "blog.codinghorror.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Envato Tuts+ Code",
				ID:      8,
				Url:     "https://code.tutsplus.com",
				FeedUrl: "https://code.tutsplus.com/posts.atom",
				Icon:    "tutsplus.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Codrops",
				ID:      9,
				Url:     "http://tympanus.net/codrops",
				FeedUrl: "http://tympanus.net/codrops/feed",
				Icon:    "codrops.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Product Hunt",
				ID:      11,
				Url:     "https://producthunt.com",
				FeedUrl: "https://posts.producthunt.com/posts/?filter=popular",
				Icon:    "producthunt.com.png",
				Crawler: "producthunt",
				Posts:   []*Post{},
			},
			{
				Title:   "Hacker News",
				ID:      12,
				Url:     "https://news.ycombinator.com",
				FeedUrl: "https://news.ycombinator.com/rss",
				Icon:    "news.ycombinator.com.ico",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "CSS Trick",
				Url:     "https://css-tricks.com/",
				ID:      13,
				FeedUrl: "http://feeds.feedburner.com/CssTricks?fmt=xml",
				Icon:    "css-tricks.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Hacker Noon",
				ID:      14,
				Url:     "https://hackernoon.com",
				FeedUrl: "https://hackernoon.com/feed",
				Icon:    "hackernoon.com.jpeg",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Tech on Medium",
				ID:      15,
				Url:     "https://medium.com/collections/d39cf943f634",
				FeedUrl: "https://medium.com/feed/collections/d39cf943f634",
				Icon:    "medium.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "TutorialZine",
				ID:      16,
				Url:     "http://tutorialzine.com",
				FeedUrl: "http://tutorialzine.com/feed/",
				Icon:    "tutorialzine.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Chromium Blog",
				ID:      17,
				Url:     "https://blog.chromium.org",
				FeedUrl: "https://blog.chromium.org/feeds/posts/default?alt=rss",
				Icon:    "blog.chromium.org.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Top story on Medium",
				ID:      18,
				Url:     "https://medium.com/browse/top",
				FeedUrl: "https://medium.com/feed/browse/top",
				Icon:    "medium.com.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Scotch.io",
				ID:      19,
				Url:     "https://scotch.io",
				FeedUrl: "https://scotch.io/feed",
				Icon:    "scotch.io.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
			{
				Title:   "Alligator",
				ID:      20,
				Url:     "https://alligator.io",
				FeedUrl: "https://alligator.io/feed.xml",
				Icon:    "alligator.io.png",
				Crawler: "rss",
				Posts:   []*Post{},
			},
		},
	}
}

func (sv *Server) Stop() {
	fmt.Println("Server stoped")
}

func (sv *Server) Start() {
	listingChan := make(chan error)
	ticker := time.NewTicker(sv.Config.CrawlInterval)
	crawlChan := make(chan string)

	go sv.Listening(listingChan)
	go sv.Crawling(crawlChan)

	for {
		select {
		case t := <-ticker.C:
			fmt.Println("Start crawling at: ", t)
			go sv.Crawling(crawlChan)

		case crawlResult := <-crawlChan:
			fmt.Println(crawlResult)

		case listingResult := <-listingChan:
			fmt.Println(listingResult)
			ticker.Stop()
			return
		}
	}
}

func (sv *Server) Listening(listingChan chan error) {

	fmt.Println("Server is listen on ", sv.Config.Address)
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/sites", sv.HandleGetSites) // Get all Sites and it's posts
	gzipWrapper := gziphandler.GzipHandler(router)

	srv := &http.Server{
		Addr:    sv.Config.Address,
		Handler: gzipWrapper,
	}
	http2.ConfigureServer(srv, nil)
	listingChan <- srv.ListenAndServe()
}

func (sv *Server) HandleGetSites(w http.ResponseWriter, r *http.Request) {
	if err := jsonapi.MarshalManyPayload(w, sv.Sites); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sv *Server) Crawling(resultChan chan string) {
	siteChan := make(chan interface{})
	chanCount := 0
	for _, site := range sv.Sites {
		chanCount++
		go sv.crawSite(site, siteChan)
	}
	for i := 0; i < chanCount; i++ {
		<-siteChan
	}
	resultChan <- "Crawl Done!"
}

func (sv *Server) crawSite(site *Site, resultChan chan interface{}) {
	var crawler CrawlerInterface.Crawler
	switch site.Crawler {
	case "rss":
		crawler = CrawlerInterface.Crawler(rss.NewCrawler(site.FeedUrl))
	case "github":
		crawler = CrawlerInterface.Crawler(github.NewCrawler(site.FeedUrl))
	case "medium":
		crawler = CrawlerInterface.Crawler(medium.NewCrawler(site.FeedUrl))
	case "producthunt":
		crawler = CrawlerInterface.Crawler(producthunt.NewCrawler(site.FeedUrl))
	case "oxfordlearnersdictionaries":
		crawler = CrawlerInterface.Crawler(oxfordlearnersdictionaries.NewCrawler(site.FeedUrl))
	}
	fmt.Println("Start parse: ", crawler.GetIdentifyURL())
	feed, err := crawler.Parse()
	if err != nil {
		resultChan <- err
		return
	}
	site.Posts = nil
	for _, item := range feed.Items {
		if len(site.Posts) >= 5 {
			break
		}
		site.Posts = append(site.Posts, &Post{
			ID:    fmt.Sprintf("%d-%d", site.ID, len(site.Posts)),
			Title: item.Title,
			Url:   item.Link,
		})
	}
	fmt.Println("End parse: ", crawler.GetIdentifyURL())
	resultChan <- true
}

// Main function
func main() {
	// UP
	configData, err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatal(err)
	}
	sv := NewServer(configData)
	defer sv.Stop()
	sv.Start()
}
