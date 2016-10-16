package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"fmt"
	"github.com/nguyenvanduocit/myfive-crawler/interface"
	"github.com/nguyenvanduocit/myfive-service/config"
	"github.com/mmcdole/gofeed"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/github"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/rss"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/medium"
	"time"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/producthunt"
	"github.com/nguyenvanduocit/myfive-crawler/crawler/oxfordlearnersdictionaries"
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

func (crawler *Crawler)insertPost(post *gofeed.Item, siteId int)(error){

	isExists, err := crawler.isPostExist(post, siteId)
	if err != nil {
		return fmt.Errorf("Insert error: %s", err)
	}
	if(isExists){
		return fmt.Errorf("Post exist: %s", post.Title)
	}
	insPost, err := crawler.db.Prepare("INSERT INTO `posts` (`site_id`, `title`, `url`, `order`) VALUES(?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		return err
	}

	result, err:= insPost.Exec(siteId , post.Title, post.Link, time.Now().UnixNano())
	if err != nil {
		return err
	}

	if _,err := result.LastInsertId(); err != nil {
		return  fmt.Errorf("Can not get post id: %s", err)
	}
	return nil
}

func (crawler *Crawler)isPostExist(post *gofeed.Item, siteId int)(bool, error){
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

func (crawler *Crawler)crawSite(crawlerClient CrawlerInterface.Crawler, resultChan chan interface{}){
	fmt.Println("Start parse: ", crawlerClient.GetIdentifyURL())
	feed , err:= crawlerClient.Parse()
	if err != nil {
		resultChan <- err
		return
	}

	siteId, err := crawler.getSiteId(crawlerClient.GetIdentifyURL())
	if err != nil {
		resultChan <- err
		return
	}
	endSlide := 5
	if(len(feed.Items) < 5){
		endSlide = len(feed.Items)
	}

	posts:= feed.Items[:endSlide]

	for i := len(posts)-1; i >= 0; i-- {
		err := crawler.insertPost(posts[i], siteId)
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println("Inserted: ", posts[i].Title)
		}
	}

	resultChan <- fmt.Sprintf("Done: %s", crawlerClient.GetIdentifyURL())
	return
}

func (crawler *Crawler)Start(){

	siteChan := make(chan interface{})

	getSiteStatement, err := crawler.db.Prepare("SELECT c.`url`,  c.`crawler` FROM `sites` as c")
	if err != nil {
		log.Fatal(err)
	}
	defer getSiteStatement.Close()

	rows, err := getSiteStatement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	chanCount := 0
	for rows.Next() {
		var url string;
		var crawlerName string;
		if err := rows.Scan(&url, &crawlerName); err != nil {
			log.Fatal(err)
		}
		switch crawlerName {
		case "rss":
			chanCount++
			go crawler.crawSite(CrawlerInterface.Crawler(rss.NewCrawler(url)), siteChan)
		case "github":
			chanCount++
			go crawler.crawSite(CrawlerInterface.Crawler(github.NewCrawler(url)), siteChan)
		case "medium":
			chanCount++
			go crawler.crawSite(CrawlerInterface.Crawler(medium.NewCrawler(url)), siteChan)
		case "producthunt":
			chanCount++
			go crawler.crawSite(CrawlerInterface.Crawler(producthunt.NewCrawler(url)), siteChan)
		case "oxfordlearnersdictionaries":
			chanCount++
			go crawler.crawSite(CrawlerInterface.Crawler(oxfordlearnersdictionaries.NewCrawler(url)), siteChan)
		}
	}
	for i := 0; i < chanCount; i++ {
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
