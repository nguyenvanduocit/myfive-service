package server

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/NYTimes/gziphandler"
	"log"
	"encoding/json"
	"html/template"
)

type Server struct{
	db *sql.DB
	ip string
	port string
}

func NewServer(dbScheme string, ip string, port string)(*Server){
	db, err := sql.Open("mysql", dbScheme)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return &Server{
		db,
		ip,
		port,
	}
}

func (sv *Server)Listing(){

	address := fmt.Sprintf("%s:%s", sv.ip, sv.port)
	fmt.Println("Server is listen on ", address);
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/sites", sv.HandleGetSites)
	router.HandleFunc("/", sv.Index)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("view/static"))))
	gzipWrapper := gziphandler.GzipHandler(router)
	log.Fatal(http.ListenAndServe(address, gzipWrapper))
}

func (sv *Server)Stop(){
	sv.db.Close()
}

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

func (sv *Server)getSites()([]*Site, error){
	getSiteStatement, err := sv.db.Prepare("SELECT c.`id`, c.`url`,  c.`title`, c.`lastupdated` FROM `sites` as c")
	if err != nil {
		return nil, err
	}
	defer getSiteStatement.Close()

	rows, err := getSiteStatement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []*Site
	for rows.Next() {
		var site Site;
		if err := rows.Scan(&site.Id, &site.Url, &site.Title, &site.LastUpdated); err != nil {
			return nil, err
		}
		site.Posts, _ = sv.getPosts(site.Id)
		sites = append(sites, &site);
	}
	return sites, nil

}

func (sv *Server)getPosts(siteId int)([]*Post, error){
	getPostsStatement, err := sv.db.Prepare("SELECT `p`.`id`, `p`.`title`, `p`.`url` FROM `posts` as `p` WHERE `p`.site_id = ? ORDER BY `p`.`pub_date` DESC LIMIT 0,5")
	if err != nil {
		return nil, err
	}
	defer getPostsStatement.Close()
	rows, err := getPostsStatement.Query(siteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post;
		if err := rows.Scan(&post.Id, &post.Title, &post.Url); err != nil {
			return nil, err
		}
		posts = append(posts, &post);
	}
	return posts, nil
}

func (sv *Server)HandleGetSites(w http.ResponseWriter, r *http.Request){
	response := &Response{
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

func (sv *Server)SendResponse(w http.ResponseWriter, r *http.Request, response *Response) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Cache-Control", "private; max-age=86400")
	json.NewEncoder(w).Encode(response)
	return
}
