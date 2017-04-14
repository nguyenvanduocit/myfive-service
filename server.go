package main

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/NYTimes/gziphandler"
	"log"
	"html/template"
	"golang.org/x/net/http2"

	"github.com/nguyenvanduocit/myfive-service/config"
	"github.com/google/jsonapi"
	"time"
	"github.com/nguyenvanduocit/myfive-crawler/interface"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/medium"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/producthunt"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/oxfordlearnersdictionaries"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/github"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/rss"
	"github.com/satori/go.uuid"
)

type Post struct {
	ID string `jsonapi:"primary,news"`
	Title string `jsonapi:"attr,title"`
	Url string `jsonapi:"attr,url"`
}

type Site struct {
	ID string `jsonapi:"primary,sites"`
	Title string `jsonapi:"attr,title"`
	Icon string `jsonapi:"attr,icon"`
	Url string `jsonapi:"attr,url"`
	FeedUrl string `jsonapi:"attr,feed_url,omitempty"`
	Crawler string
	Posts []*Post `jsonapi:"relation,posts"`
}

type Server struct{
	Config *config.Config
	Sites []*Site
}

func NewServer(config *config.Config)(*Server){
	return &Server{
		Config:config,
		Sites: []*Site{
			{
				Title: "A List Apart",
				Icon: "alistapart.com.svg",
				Url: "http://alistapart.com",
				ID: "http://alistapart.com",
				FeedUrl: "http://alistapart.com/main/feed",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Toptal",
				Icon: "toptal.com.png",
				Url: "https://www.toptal.com",
				ID: "https://www.toptal.com",
				FeedUrl: "https://www.toptal.com/blog.rss",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Smashing Magazine",
				Url: "https://www.smashingmagazine.com",
				ID: "https://www.smashingmagazine.com",
				FeedUrl: "https://www.smashingmagazine.com/feed/",
				Icon: "toptal.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "David Walsh Blog",
				Url: "https://davidwalsh.name",
				ID: "https://davidwalsh.name",
				FeedUrl: "https://davidwalsh.name/feed",
				Icon: "davidwalsh.name.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "SitePoint",
				Url: "https://www.sitepoint.com",
				ID: "https://www.sitepoint.com",
				FeedUrl: "https://www.sitepoint.com/feed",
				Icon: "sitepoint.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "GitHub Trending",
				ID: "https://github.com/trending",
				Url: "https://github.com/trending",
				FeedUrl: "https://github.com/trending",
				Icon: "github.svg",
				Crawler: "github",
				Posts: []*Post{},
			},
			{
				Title: "Echo JS",
				ID: "http://www.echojs.com",
				Url: "http://www.echojs.com",
				FeedUrl: "http://www.echojs.com/rss",
				Icon: "echojs.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Coding Horror",
				ID: "https://blog.codinghorror.com",
				Url: "https://blog.codinghorror.com",
				FeedUrl: "https://blog.codinghorror.com/rss/",
				Icon: "blog.codinghorror.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Envato Tuts+ Code",
				ID: "https://code.tutsplus.com",
				Url: "https://code.tutsplus.com",
				FeedUrl: "https://code.tutsplus.com/posts.atom",
				Icon: "tutsplus.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Codrops",
				ID: "http://tympanus.net/codrops",
				Url: "http://tympanus.net/codrops",
				FeedUrl: "http://tympanus.net/codrops/feed",
				Icon: "codrops.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "WordPress Tavern",
				ID: "https://wptavern.com",
				Url: "https://wptavern.com",
				FeedUrl: "https://wptavern.com/feed",
				Icon: "wptavern.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Product Hunt",
				ID: "https://producthunt.com",
				Url: "https://producthunt.com",
				FeedUrl: "https://posts.producthunt.com/posts/?filter=popular",
				Icon: "producthunt.com.png",
				Crawler: "producthunt",
				Posts: []*Post{},
			},
			{
				Title: "Hacker News",
				ID: "https://news.ycombinator.com",
				Url: "https://news.ycombinator.com",
				FeedUrl: "https://news.ycombinator.com/rss",
				Icon: "news.ycombinator.com.ico",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "CSS Trick",
				Url: "https://css-tricks.com/",
				FeedUrl: "http://feeds.feedburner.com/CssTricks?fmt=xml",
				Icon: "css-tricks.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Hacker Noon",
				ID: "https://hackernoon.com",
				Url: "https://hackernoon.com",
				FeedUrl: "https://hackernoon.com/feed",
				Icon: "hackernoon.com.jpeg",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Tech on Medium",
				ID: "https://medium.com/collections/d39cf943f634",
				Url: "https://medium.com/collections/d39cf943f634",
				FeedUrl: "https://medium.com/feed/collections/d39cf943f634",
				Icon: "medium.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "TutorialZine",
				ID: "http://tutorialzine.com",
				Url: "http://tutorialzine.com",
				FeedUrl: "http://tutorialzine.com/feed/",
				Icon: "tutorialzine.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Chromium Blog",
				ID: "https://blog.chromium.org",
				Url: "https://blog.chromium.org",
				FeedUrl: "https://blog.chromium.org/feeds/posts/default?alt=rss",
				Icon: "blog.chromium.org.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Top story on Medium",
				ID: "https://medium.com/browse/top",
				Url: "https://medium.com/browse/top",
				FeedUrl: "https://medium.com/feed/browse/top",
				Icon: "medium.com.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Scotch.io",
				ID: "https://scotch.io",
				Url: "https://scotch.io",
				FeedUrl: "https://scotch.io/feed",
				Icon: "scotch.io.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
			{
				Title: "Alligator.io",
				ID: "https://alligator.io",
				Url: "https://alligator.io",
				FeedUrl: "https://alligator.io/feed.xml",
				Icon: "alligator.io.png",
				Crawler: "rss",
				Posts: []*Post{},
			},
		},
	}
}

func (sv *Server)Stop(){
	fmt.Println("Server stoped");
}

func (sv *Server)Start(){
	listingChan := make(chan error)
	ticker := time.NewTicker(time.Minute)
	crawlChan := make(chan string)

	go sv.Listing(listingChan)
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

func (sv *Server)Listing(listingChan chan error){

	fmt.Println("Server is listen on ", sv.Config.Address);
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/sites", sv.HandleGetSites) // Get all Sites and it's posts
	//router.HandleFunc("/api/v1/picked_news", sv.HandleGetPickedNews) // Get 5 picks by developer
	//router.HandleFunc("/api/v1/pick_news", sv.HandlePickNews) // Handle post new pick
	router.HandleFunc("/", sv.Index)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("view/static"))))
	gzipWrapper := gziphandler.GzipHandler(router)

	srv := &http.Server{
		Addr:    sv.Config.Address,
		Handler: gzipWrapper,
	}
	http2.ConfigureServer(srv, nil)
	listingChan <- srv.ListenAndServe()
}

func (sv *Server)HandleGetSites(w http.ResponseWriter, r *http.Request){
	if err := jsonapi.MarshalManyPayload(w, sv.Sites); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sv *Server)Crawling(resultChan chan string){
	siteChan := make(chan interface{})
	chanCount := 0
	for _, site  := range sv.Sites {
		chanCount++
		go sv.crawSite( site, siteChan)
	}
	for i := 0; i < chanCount; i++ {
		<-siteChan
	}
	resultChan <- "Crawl Done!"
}

func (sv *Server)crawSite(site *Site, resultChan chan interface{}){
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
	fmt.Println( "Start parse: ", crawler.GetIdentifyURL())
	feed , err:= crawler.Parse()
	if err != nil {
		resultChan <- err
		return
	}
	site.Posts = nil
	for _,item := range feed.Items {
		if len(site.Posts) >= 5 {
			break
		}
		site.Posts = append(site.Posts, &Post{
			ID: uuid.NewV4().String(),
			Title: item.Title,
			Url: item.Link,
		})
	}
	fmt.Println("End parse: ", crawler.GetIdentifyURL())
	resultChan <- true
}
// Handle index request
func (sv *Server)Index(w http.ResponseWriter, r *http.Request){
	templates, err := template.ParseFiles( "./view/index.html" );
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = templates.Execute(w,  nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Main function

func main() {
	// UP
	configData, err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatal(err)
	}
	sv := NewServer(configData);
	defer sv.Stop();
	sv.Start();
}
