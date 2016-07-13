package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"github.com/nguyenvanduocit/myfive-service/server"
)

type Installer struct {
	db *sql.DB
}

func NewInstaller(dbScheme string)(*Installer){
	db, err := sql.Open("mysql", dbScheme)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return &Installer{
		db,
	}
}

func (installer *Installer)importPosts(siteId int64, posts []*server.Post){
	insPost, err := installer.db.Prepare("INSERT INTO `posts` (site_id, title, url) VALUES(?, ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < len(posts); i++ {
		result, err := insPost.Exec(siteId , posts[i].Title, posts[i].Url)
		if err != nil {
			panic(err.Error())
		}
		postId, err := result.LastInsertId()
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Post:", postId)
	}
}

func (installer *Installer)ImportSite(fileName string){

	siteFile, err := os.Open("data/" + fileName)
	if err != nil {
		panic(err.Error())
	}
	defer siteFile.Close()
	var site server.Site
	jsonParser := json.NewDecoder(siteFile)
	err = jsonParser.Decode(&site);
	if err != nil {
		panic(err.Error())
	}

	insSite, err := installer.db.Prepare("INSERT INTO `sites` (url, title) VALUES(?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}
	result, err := insSite.Exec(site.Url, site.Title)
	if err != nil {
		panic(err.Error())
	}
	siteId, err := result.LastInsertId();
	if err != nil {
		panic(err.Error())
	}
	//TODO Use chan
	installer.importPosts(siteId, site.Posts)
}

func (installer *Installer)CreateTable(name string,query string){

	if _, err := installer.db.Exec("SET foreign_key_checks = 0"); err != nil {
		panic(err.Error())
	}

	if _, err := installer.db.Exec("DROP TABLE IF EXISTS `"+ name + "`"); err != nil {
		panic(err.Error())
	}

	if _, err := installer.db.Exec(query); err != nil {
		panic(err.Error())
	}
	fmt.Println("Table " + name + " created")
}

func (installer *Installer)Done(){
	installer.db.Close()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbScheme := os.Getenv("db-scheme")
	installer:= NewInstaller(dbScheme)
	defer installer.Done()

	installer.CreateTable("sites", "CREATE TABLE `sites` ( `id` int(11) unsigned NOT NULL AUTO_INCREMENT, `title` varchar(255) DEFAULT NULL, `url` varchar(255) DEFAULT NULL, `lastupdated` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	installer.CreateTable("posts", "CREATE TABLE `posts` ( `id` int(11) unsigned NOT NULL AUTO_INCREMENT, `site_id` int(11) unsigned, `title` text,`url` varchar(255) DEFAULT NULL, `pub_date` timestamp DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), FOREIGN KEY (site_id) REFERENCES sites(id) ) ENGINE=InnoDB AUTO_INCREMENT=424 DEFAULT CHARSET=utf8;")

	siteFiles, _ := ioutil.ReadDir("./data")
	for _, f := range siteFiles {
		fmt.Println("Import : " + f.Name())
		installer.ImportSite(f.Name())
	}
}
