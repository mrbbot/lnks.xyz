package main

import (
	_ "github.com/joho/godotenv/autoload"
	_ "shortener/db"
	"shortener/routes"
)

func main() {
	n := routes.NewRouter()
	n.Run()
}
