package main

import (
	"log"
	"net/http"

	"main.go/Connection"
)

const (
	Host = "localhost"
	Port = "2077"
)

func main() {
	Connection.Connection()
	err := http.ListenAndServe(Host+":"+Port, nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
		return
	}
}
