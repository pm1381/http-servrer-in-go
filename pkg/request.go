package pkg

type UploadRequest struct {
	File string `json:"file"`
}

type DownloadRequest struct {
	File_id string `json:"file_id"`
}