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
	"github.com/nguyenvanduocit/myfive-service/database"
	"github.com/google/jsonapi"
)

type Server struct{
	Config *config.Config
	DbFactory *database.DbFactory
}

func NewServer(config *config.Config)(*Server){
	return &Server{
		Config:config,
		DbFactory: &database.DbFactory{
			DatabaseName:config.DatabaseName,
			Host:config.DatabaseHost,
			Port:config.DatabasePort,
			Username:config.DatabaseUserName,
			Password:config.DatabasePassword,
		},
	}
}

func (sv *Server)Stop(){
	fmt.Println("Server stoped");
}
func (sv *Server)Listing(){

	fmt.Println("Server is listen on ", sv.Config.Address);
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/sites", sv.HandleGetSites)
	router.HandleFunc("/", sv.Index)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("view/static"))))
	gzipWrapper := gziphandler.GzipHandler(router)

	srv := &http.Server{
		Addr:    sv.Config.Address,
		Handler: gzipWrapper,
	}
	http2.ConfigureServer(srv, nil)
	log.Panic(srv.ListenAndServe())
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

// Handle get Sites
func (sv *Server)HandleGetSites(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", jsonapi.MediaType)
	sites, err := sv.getSites()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := jsonapi.MarshalManyPayload(w, sites); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sv *Server)getSites()([]*database.Site, error){
	db := sv.DbFactory.NewConnect()
	defer db.Close()
	getSiteStatement, err := db.Prepare("SELECT c.`id`, c.`url`, c.`icon`,  c.`title` FROM `sites` as c")
	if err != nil {
		return nil, err
	}
	defer getSiteStatement.Close()

	rows, err := getSiteStatement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []*database.Site
	for rows.Next() {
		var site database.Site;
		if err := rows.Scan(&site.Id, &site.Url, &site.Icon, &site.Title); err != nil {
			return nil, err
		}
		site.Posts, _ = sv.getNewsList(site.Id)
		sites = append(sites, &site);
	}
	return sites, nil

}

func (sv *Server)getNewsList(siteId int)([]*database.News, error){
	db := sv.DbFactory.NewConnect()
	defer db.Close()
	getNewsListStatement, err := db.Prepare("SELECT `p`.`id`, `p`.`title`, `p`.`url` FROM `posts` as `p` WHERE `p`.site_id = ? ORDER BY `p`.`id` DESC LIMIT 0,5")
	if err != nil {
		return nil, err
	}
	defer getNewsListStatement.Close()
	rows, err := getNewsListStatement.Query(siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []*database.News
	for rows.Next() {
		var news database.News;
		if err := rows.Scan(&news.Id, &news.Title, &news.Url); err != nil {
			return nil, err
		}
		newsList = append(newsList, &news);
	}
	return newsList, nil
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
	sv.Listing();
}
