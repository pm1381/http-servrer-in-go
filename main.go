package main

import (
	"intern/http-server/internal/http_server"

	"github.com/gorilla/mux"
)

func main()  {
	router := mux.NewRouter()
	http_server.Run(router)
}