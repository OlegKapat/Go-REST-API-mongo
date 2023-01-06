package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/OlegKapat/Rest-api-mongo/internal/user"
	"github.com/julienschmidt/httprouter"
)

func main() {
	log.Println("Create router")
	router := httprouter.New()

	log.Println("Register user handler")
	handler := user.NewHandler()
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
	log.Println("Start application")
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	server := http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server is listering port 0.0.0.0:1234")
	log.Fatalln(server.Serve(listener))
}
