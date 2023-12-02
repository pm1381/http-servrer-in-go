package api

import (
	"encoding/json"
	"fmt"
	"intern/http-server/pkg"
	"net/http"
)

func JsonUpload(w http.ResponseWriter, r *http.Request)  {
	var request pkg.UploadRequest
	iss := ""
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {iss = "error in user params"}
	
	filePath := pkg.ReadFromUrl(request.File)
	if filePath == "" {iss = "no file found or created"}
	
	file := pkg.OpenFile(filePath)
	if (file == nil) {iss = "cannot open source file"}
	defer file.Close()

	destFile := pkg.WriteToFile(file.Name())
	if destFile == nil {iss = "destination file not found"}
	defer destFile.Close()
	if iss == "" {
		iss = pkg.UpProcess(file, destFile)
	}
	if iss != "" {
		mapData := map[string]interface{}{
			"error" : iss,
		}
		jsonData, _ := json.Marshal(mapData)
		pkg.JsonFormat(w, http.StatusInternalServerError, jsonData)
	} else {
		mapData := map[string]interface{}{
			"file_id" : destFile.Name(),
		}
		jsonData, _ := json.Marshal(mapData)
		pkg.JsonFormat(w, http.StatusInternalServerError, jsonData)
	}
}

func FormUpload(w http.ResponseWriter, r *http.Request)  {
	iss := ""
	file, header, err := r.FormFile("file")
	if err != nil {iss = "error in user params"}
	defer file.Close()
	ioFile := pkg.ReadFromForm(header, file)
	if (ioFile == nil) {iss = "source file not found"}
	ioFile.Close()
	destFile := pkg.WriteToFile(ioFile.Name())
	if (destFile == nil) {iss = "destination file not found"}
	ioFile.Close()
	if iss == "" {
		iss = pkg.UpProcess(ioFile, destFile)
	}
	if iss != "" {
		mapData := map[string]interface{}{
			"error" : iss,
		}
		jsonData, _ := json.Marshal(mapData)
		pkg.JsonFormat(w, http.StatusInternalServerError, jsonData)
	} else {
		mapData := map[string]interface{}{
			"file_id" : destFile.Name(),
		}
		jsonData, _ := json.Marshal(mapData)
		pkg.JsonFormat(w, http.StatusInternalServerError, jsonData)
	}
}

func FormDownload(w http.ResponseWriter, r *http.Request)  {
	iss := ""
	file_id := r.FormValue("file_id")
	if file_id == "" {
		iss = "destination file not found"
	}
	file := pkg.OpenFile("upload/" + file_id + ".txt")
	res, iss := pkg.DlProcess(file)
	if (iss == "") {
		for _, v := range res {
			fmt.Fprintf(w, string(v))
			w.WriteHeader(http.StatusOK)
		}
	} else {
		mapData := map[string]interface{}{
			"error": iss,
		}
		jsonData, _ := json.Marshal(mapData)
		pkg.JsonFormat(w, http.StatusInternalServerError, jsonData)
	}
}

func JsonDownload(w http.ResponseWriter, r *http.Request)  {
	var request pkg.DownloadRequest
	iss := ""
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {iss = "error in user params"}
	file := pkg.OpenFile("upload/" + request.File_id + ".txt")
	res, iss := pkg.DlProcess(file)
	if (iss == "") {
		for _, v := range res {
			fmt.Fprintf(w, string(v))
			w.WriteHeader(http.StatusOK)
		}
	} else {
		mapData := map[string]interface{}{
			"error": iss,
		}
		jsonData, _ := json.Marshal(mapData)
		pkg.JsonFormat(w, http.StatusInternalServerError, jsonData)
	}
}