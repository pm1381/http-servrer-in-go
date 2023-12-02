package http_server

import (
	"intern/http-server/api"
	"net/http"

	"github.com/gorilla/mux"
)

func Run(router *mux.Router)  {
	router.HandleFunc("/uploadFile", api.JsonUpload).Methods("POST").Headers("Content-Type", "application/json")
	// router.HandleFunc("/uploadFile", formUpload).Methods("POST").Headers("Content-Type", "multipart/form-data")
	router.HandleFunc("/uploadFile", api.FormUpload).Methods("POST")
	router.HandleFunc("/downloadFile", api.JsonDownload).Methods("GET").Headers("Content-Type", "application/json")
	router.HandleFunc("/downloadFile", api.FormDownload).Methods("GET")
	http.ListenAndServe(":80", router)
}