package main

import (
	"fmt"
	"github.com/NYTimes/gziphandler"
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
	"github.com/rs/cors"
	"github.com/graphql-go/handler"
	"github.com/nguyenvanduocit/myfive-service/schema"
)

type Post struct {
	ID    string `json:"id" jsonapi:"primary,news"`
	Title string `json:"title" jsonapi:"attr,title"`
	Url   string `json:"url" jsonapi:"attr,url"`
}

type Site struct {
	ID      int    `json:"id" jsonapi:"primary,sites"`
	Title   string `json:"title" jsonapi:"attr,title"`
	Icon    string `json:"icon" jsonapi:"attr,icon"`
	Url     string `json:"url" jsonapi:"attr,url"`
	FeedUrl string
	Crawler string
	Posts   []*Post `json:"posts" jsonapi:"relation,posts"`
	LastUpdated time.Time
}

type Server struct {
	Config *config.Config
	Sites  []*Site
	BatchSize int
	CacheInterval time.Duration
}

func NewServer(config *config.Config) *Server {
	return &Server{
		Config:        config,
		BatchSize:     5,
		CacheInterval: 20 * time.Minute,
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
				Title:   "freeCodeCamp",
				ID:      15,
				Url:     "https://medium.freecodecamp.com",
				FeedUrl: "https://medium.freecodecamp.com/feed",
				Icon:    "medium.com.png",
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
			},{
				Title:   "Google Developers Blog",
				Url:     "https://developers.googleblog.com/",
				ID:      22,
				FeedUrl: "https://www.blogger.com/feeds/596098824972435195/posts/default",
				Icon:    "developers.googleblog.com.jpg",
				Crawler: "rss",
				Posts:   []*Post{},
			},{
				Title:   "Bitsofco",
				ID:      23,
				Url:     "https://bitsofco.de",
				FeedUrl: "https://bitsofco.de/rss",
				Icon:    "bitsofco.de.png",
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
	crawlChan := make(chan bool)
	previousDone := false
	go sv.Listening(listingChan)
	go sv.Crawling(crawlChan)

	for {
		select {
		case <-ticker.C:
			if previousDone == true {
				previousDone = false
				go sv.Crawling(crawlChan)
			}

		case <-crawlChan:
			previousDone = true

		case listingResult := <-listingChan:
			fmt.Println(listingResult)
			ticker.Stop()
			return
		}
	}
}

func (sv *Server) getSites (id int) *Site {
	for _, site := range sv.Sites {
		if site.ID == id {
			return site
		}
	}
	return nil
}

func (sv *Server) Listening(listingChan chan error) {
	fmt.Println("Server is listen on ", sv.Config.Address)
	router := mux.NewRouter().StrictSlash(true)
	corsMiddleWare := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	router.HandleFunc("/api/v1/sites", sv.HandleGetSites) // Get all Sites and it's posts

	graphqlHandler := handler.New(&handler.Config{
		Schema: &schema.Schema,
		Pretty: true,
	})
	router.Handle("/graphql", graphqlHandler)

	gzipWrapper := gziphandler.GzipHandler(corsMiddleWare.Handler(router))

	srv := &http.Server{
		Addr:    sv.Config.Address,
		Handler: gzipWrapper,
	}
	http2.ConfigureServer(srv, nil)
	listingChan <- srv.ListenAndServe()
}

func (sv *Server) HandleGetSites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := jsonapi.MarshalManyPayload(w, sv.Sites); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sv *Server) Crawling(resultChan chan bool) {
	timeNow := time.Now()
	fmt.Printf("\n---| %s \t Start patch!\n", timeNow)
	siteChan := make(chan interface{})
	chanCount := 0
	for _, site := range sv.Sites {
		isFirstTime := site.LastUpdated.IsZero()
		if  isFirstTime || timeNow.Sub(site.LastUpdated) > sv.CacheInterval {
			chanCount++
			site.LastUpdated = timeNow
			fmt.Printf("%s \t Start parse:\t %s\n", site.LastUpdated, site.Title)
			go sv.crawSite(site, siteChan)
		}
		if (chanCount == sv.BatchSize) && !isFirstTime {
			break
		}
	}
	for i := 0; i < chanCount; i++ {
		fmt.Printf("%s \t End parse:\t\t %s\n", time.Now(), <-siteChan)
	}
	fmt.Printf("---| %s \t End patch!\n", time.Now())
	resultChan <- true
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
	resultChan <- site.Title
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
