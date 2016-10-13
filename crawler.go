package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"fmt"
	"github.com/nguyenvanduocit/myfive-crawler/interface"
	"github.com/nguyenvanduocit/myfive-crawler/crawler"
	"github.com/nguyenvanduocit/myfive-crawler/model/rss"
	"time"
	"github.com/nguyenvanduocit/myfive-crawler/model/site"
	"github.com/nguyenvanduocit/myfive-service/config"
)

type Crawler struct {
	db *sql.DB
}

func NewCrawler(sbScheme string)(*Crawler){
	db, err := sql.Open("mysql", sbScheme)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return &Crawler{
		db,
	}
}

func (crawler *Crawler)insertPost(post *RssModel.Item, siteId int, ch chan interface{}){

	isExists, err := crawler.isPostExist(post, siteId)
	if err != nil {
		ch <- fmt.Errorf("Insert error: %s", err)
		return
	}
	if(isExists){
		ch <- fmt.Errorf("Post exist: %s", post.Title)
		return
	}

	insPost, err := crawler.db.Prepare("INSERT INTO `posts` (site_id, title, url, pub_date) VALUES(?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		ch <- err
		return
	}
	pubDate, err := time.Parse("Mon, _2 Jan 2006 15:04:05 +0000",post.PubDate)
	result, err:= insPost.Exec(siteId , post.Title, post.Link, pubDate.Format("2006-01-02 15:04:05"))
	if err != nil {
		ch <- err
		return
	}

	if _,err := result.LastInsertId(); err != nil {
		ch <- fmt.Errorf("Can not get post id: %s", err)
		return
	}
	ch <- fmt.Sprintf("Insert success: %s", post.Title)
	return
}

func (crawler *Crawler)isPostExist(post *RssModel.Item, siteId int)(bool, error){
	checkStatement, err := crawler.db.Prepare("SELECT EXISTS(SELECT 1 FROM `posts` as `p` WHERE `p`.url = ?)") // ? = placeholder
	if err != nil {
		return false, err
	}
	var exists bool
	if err := checkStatement.QueryRow(post.Link).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil

}

func (crawler *Crawler)getSiteId(url string)(int, error){
	insertStatement, err := crawler.db.Prepare("SELECT `s`.`id` FROM `sites` as `s` WHERE `s`.url = ?") // ? = placeholder
	if err != nil {
		return -1, err
	}
	var siteId int

	if err := insertStatement.QueryRow(url).Scan(&siteId); err != nil {
		return -1, err
	}
	return siteId, nil

}

func (crawler *Crawler)insertSite(site *SiteModel.Site)(int, error){
	insSite, err := crawler.db.Prepare("INSERT INTO `sites` (url, title) VALUES(?, ? )") // ? = placeholder
	if err != nil {
		return -1, err
	}
	result, err := insSite.Exec(site.Link, site.Title)
	if err != nil {
		return -1, err
	}
	siteId, err := result.LastInsertId();
	if err != nil {
		return -1, err
	}
	return int(siteId), nil

}

func (crawler *Crawler)crawSite(url string, crawlerClient CrawlerInterface.Crawler, resultChan chan interface{}){
	siteId, err := crawler.getSiteId(url)
	if err != nil {
		siteInfo, err := crawlerClient.GetSiteInfo()
		if err != nil {
			resultChan <- err
			return
		}
		siteId, err = crawler.insertSite(siteInfo)
		if err != nil {
			resultChan <- err
			return
		}
	}
	posts, err := crawlerClient.GetTopFive()
	if err != nil {
		resultChan <- err
		return
	}

	postChan := make(chan interface{})

	for _,post := range posts{
		go crawler.insertPost(post, siteId, postChan)
	}
	totalPost := len(posts)
	for i:=0; i<totalPost; i++ {
		fmt.Println( <- postChan)
	}

	resultChan <- fmt.Sprintf("Done: %s", url)
	return
}

func (crawler *Crawler)Start(){

	siteChan := make(chan interface{})

	go crawler.crawSite("http://sitepoint.com", CrawlerInterface.Crawler(RssCrawler.NewCrawler("https://www.sitepoint.com/feed")), siteChan)
	go crawler.crawSite("https://davidwalsh.name", CrawlerInterface.Crawler(RssCrawler.NewCrawler("https://davidwalsh.name/feed")), siteChan)
	go crawler.crawSite("https://wptavern.com", CrawlerInterface.Crawler(RssCrawler.NewCrawler("https://wptavern.com/feed")), siteChan)
	go crawler.crawSite("https://laptrinh.senviet.org", CrawlerInterface.Crawler(RssCrawler.NewCrawler("https://laptrinh.senviet.org/feed")), siteChan)

	for i := 0; i < 4; i++ {
		fmt.Println(<-siteChan)
	}
}

func (crawler *Crawler)Done(){
	crawler.db.Close()
}

func main() {
	configData, err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatal(err)
	}
	dbScheme := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", configData.DatabaseUserName, configData.DatabasePassword, configData.DatabaseHost, configData.DatabasePort, configData.DatabaseName)
	crawler := NewCrawler(dbScheme)
	defer crawler.Done()
	crawler.Start()

}
