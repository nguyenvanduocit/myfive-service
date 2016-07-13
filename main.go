package main

import (
	"log"
)
import (
	"github.com/joho/godotenv"
	"os"
	"github.com/nguyenvanduocit/myfive-service/server"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ip := os.Getenv("ip")
	port := os.Getenv("port")
	dbScheme := os.Getenv("db-scheme")

	sv := server.NewServer(dbScheme, ip, port);
	sv.Listing();
	defer sv.Stop();
}
