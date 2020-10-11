package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// var filepath = "./"
// var filename = "android-studio-ide-193.6821437-windows.exe"

func handlerCloud(w http.ResponseWriter, req *http.Request) {
	if q := *req.URL; q.RawQuery != "" {
		// fmt.Println(q)
		filePath := q.Path[7:]
		fileName := q.RawQuery
		// fmt.Println(filepath)
		// fmt.Println(filename)
		if filePath == "" {
			filePath = "."
		}

		file, err := os.Open(filePath + "/" + fileName)
		if err != nil {
			log.Println(err)
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileHeader := make([]byte, 512)
		file.Read(fileHeader)
		fileStat, _ := file.Stat()

		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, fileName))
		w.Header().Set("Content-Type", http.DetectContentType(fileHeader))

		var start, end int64
		if r := req.Header.Get("Range"); r != "" {
			fmt.Println(r)
			if strings.Contains(r, "bytes=") && strings.Contains(r, "-") {
				fmt.Sscanf(r, "bytes=%d-%d", &start, &end)
				if end == 0 {
					end = fileStat.Size() - 1
				}
				// fmt.Println(start, end)
				w.Header().Set("Content-Length", fmt.Sprintf("%v", end-start+1))
				w.Header().Set("Content-Range", fmt.Sprintf("bytes %v-%v/%v", start, end, fileStat.Size()))
				w.WriteHeader(http.StatusPartialContent)
			}
		} else {
			w.Header().Set("Content-Length", fmt.Sprintf("%v", fileStat.Size()))
			start = 0
			end = fileStat.Size() - 1
		}

		file.Seek(start, 0)
		var n int64 = 512
		buf := make([]byte, n)
		for start <= end {
			if end-start+1 < n {
				n = end - start + 1
			}
			file.Read(buf[:n])
			w.Write(buf[:n])
			start += int64(n)
		}

	} else {

	}
}

func main() {
	http.HandleFunc("/cloud/", handlerCloud)
	http.ListenAndServe(":9999", nil)
}
