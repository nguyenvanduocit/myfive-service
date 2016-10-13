package database

type Post struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Url string `json:"url"`
}

type Site struct{
	Id int `json:"id"`
	Url string `json:"url"`
	Title string `json:"title"`
	LastUpdated string `json:"lastupdated"`
	Posts []*Post `json:"posts"`
}

type Response struct{
	Success bool `json:"success"`
	Message string `json:"message"`
	Count int `json:"count"`
	Result interface{} `json:"result"`
}
