package pkg

import "net/http"

	
func SimpleFormat(w http.ResponseWriter, code int)  {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
}

func JsonFormat(w http.ResponseWriter, code int, jsonData []byte)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)
}