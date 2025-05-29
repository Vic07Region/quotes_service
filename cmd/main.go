package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"quotes_service/internal/server"
	"strconv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	addrEnv := os.Getenv("QS_RUN_HOST")
	if addrEnv == "" {
		addrEnv = ":8080"
	}
	debugEnv := os.Getenv("QS_DEBUG")
	d, err := strconv.ParseBool(debugEnv)
	if err != nil {
		d = false
	}
	srv := server.NewServer(addrEnv, d)
	srv.StartServer()
}
