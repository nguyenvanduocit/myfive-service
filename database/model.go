package database

type News struct {
	Id int `jsonapi:"primary,news"`
	Title string `jsonapi:"attr,title"`
	Url string `jsonapi:"attr,url"`
}

type Site struct{
	Id int `jsonapi:"primary,sites"`
	Url string `jsonapi:"attr,url"`
	Icon string `jsonapi:"attr,icon"`
	FeedUrl string `jsonapi:"attr,feed_url,omitempty"`
	Title string `jsonapi:"attr,title"`
	Posts []*News `jsonapi:"relation,news"`
}
