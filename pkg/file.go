package pkg

import (
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

const bufferSize = 100

func UpProcess(file , destFile  *os.File) (string)  {
	concurrencySteps, fileSize := findNumberOfSteps(file, bufferSize)
	if concurrencySteps == 0 || fileSize == 0 {
		return "file is empty"
	}
	chunkArray := fillChunkSizes(concurrencySteps, fileSize)
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(concurrencySteps)
	for i := 0; i < concurrencySteps; i++ {
		go func (chunks []chunk, i int, mu *sync.Mutex, wg *sync.WaitGroup,destFile *os.File)  {
			defer wg.Done()
			mu.Lock()
			chunk := chunks[i]
			buffer := make([]byte, chunk.bufSize)
			file.ReadAt(buffer, chunk.offset);
			destFile.WriteAt(buffer, chunk.offset)
			mu.Unlock()
		}(chunkArray, i, &mu, &wg, destFile)
	}
	wg.Wait()
	return ""
}

func DlProcess(file  *os.File) ([][]byte, string)  {
	concurrencySteps, fileSize := findNumberOfSteps(file, bufferSize)
	data := make([][]byte, concurrencySteps)
	if concurrencySteps == 0 || fileSize == 0 {
		return data, "file is empty"
	}
	chunkArray := fillChunkSizes(concurrencySteps, fileSize)
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(concurrencySteps)
	for i := 0; i < concurrencySteps; i++ {
		go func (chunks []chunk, i int, mu *sync.Mutex, wg *sync.WaitGroup, data *[][]byte)  {
			defer wg.Done()
			mu.Lock()
			chunk := chunks[i]
			buffer := make([]byte, chunk.bufSize)
			bytesRead, _ := file.ReadAt(buffer, chunk.offset);
			(*data)[i] = ((buffer[:bytesRead]))
			mu.Unlock()
		}(chunkArray, i, &mu, &wg, &data)
	}
	wg.Wait()
	return data, ""
}

func findNumberOfSteps(file *os.File, bufferSize int) (int,int) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, 0
	}
	fileSize := int(fileInfo.Size())
	concurrencySteps := fileSize / bufferSize
	if fileSize % bufferSize != 0 {
		concurrencySteps++
	}
	if concurrencySteps == 0 || bufferSize == 0 {
		return 0, 0
	}
	return concurrencySteps, fileSize
}

func OpenFile(fileName string) (*os.File) {
	file, err := os.Open(fileName)
	if err != nil {return nil}
	return file
}

func ReadFromUrl(url  string) (string)  {
	resp, err := http.Get(url)
	if err != nil {return ""}
	defer resp.Body.Close()
	filepath := "temp/" + strings.ReplaceAll(path.Base(resp.Request.URL.String()), ".", "") + ".txt"
	out, err := os.Create(filepath)
	if err != nil {return ""}
	defer out.Close()
	_, _ = io.Copy(out, resp.Body)
	return filepath
}

func ReadFromForm(header *multipart.FileHeader, file multipart.File) (*os.File) {
	filepath := "temp/" + header.Filename
	outputFile, err := os.Create(filepath)
	if err != nil {
		return nil
	}
	defer outputFile.Close()
	_, _ = io.Copy(outputFile, file)
	return outputFile
}

func WriteToFile(initName string) (*os.File) {
	updatedName := strconv.FormatUint(hash(initName), 10) + "_" + base64.StdEncoding.EncodeToString([]byte(initName)) +	".txt"
	updatedName = strings.ReplaceAll(updatedName, "/", "_")
	updatedName = "upload/" + updatedName
	file, err := os.OpenFile(updatedName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {return nil}
	return file
}
