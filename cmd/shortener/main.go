package main

import (
	"github.com/joho/godotenv"
	"log"
	"ozon_entrance/internal/server"
)

func main() {
	_ = godotenv.Load()
	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
