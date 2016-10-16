package main

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/NYTimes/gziphandler"
	"log"
	"encoding/json"
	"html/template"
	"golang.org/x/net/http2"

	"github.com/nguyenvanduocit/myfive-service/config"
	"github.com/nguyenvanduocit/myfive-service/database"
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

	router.HandleFunc("/api/v1/news", sv.HandleGetSites)
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

// Handler functions

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

func (sv *Server)HandleGetSites(w http.ResponseWriter, r *http.Request){
	response := &database.Response{
		Success:false,
		Message:"Unknown error!",
	}
	sites, err := sv.getSites()
	if (err != nil) {
		response.Message = err.Error()
	}else{
		response.Message= "Success"
		response.Success = true
		response.Result = sites
		response.Count = len(sites)
	}
	sv.SendResponse(w, r, response)
	return
}

// Until functions

func (sv *Server)SendResponse(w http.ResponseWriter, r *http.Request, response *database.Response) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Cache-Control", "private; max-age=86400")
	json.NewEncoder(w).Encode(response)
	return
}

// Repository functions

func (sv *Server)getSites()([]*database.Site, error){
	db := sv.DbFactory.NewConnect()
	defer db.Close()
	getSiteStatement, err := db.Prepare("SELECT c.`id`, c.`url`,  c.`title` FROM `sites` as c")
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
		if err := rows.Scan(&site.Id, &site.Url, &site.Title); err != nil {
			return nil, err
		}
		site.Posts, _ = sv.getPosts(site.Id)
		sites = append(sites, &site);
	}
	return sites, nil

}

func (sv *Server)getPosts(siteId int)([]*database.Post, error){
	db := sv.DbFactory.NewConnect()
	defer db.Close()
	getPostsStatement, err := db.Prepare("SELECT `p`.`id`, `p`.`title`, `p`.`url` FROM `posts` as `p` WHERE `p`.site_id = ? ORDER BY `p`.`order` DESC LIMIT 0,5")
	if err != nil {
		return nil, err
	}
	defer getPostsStatement.Close()
	rows, err := getPostsStatement.Query(siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*database.Post
	for rows.Next() {
		var post database.Post;
		if err := rows.Scan(&post.Id, &post.Title, &post.Url); err != nil {
			return nil, err
		}
		posts = append(posts, &post);
	}
	return posts, nil
}

// Main function

func main() {
	// UP
	configData, err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatal(err)
	}
	sv := NewServer(configData);
	sv.Listing();
	defer sv.Stop();
}
