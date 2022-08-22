package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	diff_time = 3600 * 24 * 7
)

// 获取文件大小的接口
type Size interface {
	Size() int64
}

// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}

func uploadServer(w http.ResponseWriter, r *http.Request) {

	if "POST" == r.Method {
		file, _, err := r.FormFile("userfile")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if statInterface, ok := file.(Stat); ok {
			fileInfo, _ := statInterface.Stat()
			fmt.Fprintf(w, "上传文件的大小为: %d", fileInfo.Size())
		}
		if sizeInterface, ok := file.(Size); ok {
			fmt.Fprintf(w, "上传文件的大小为: %d", sizeInterface.Size())
		}
		err2 := r.ParseMultipartForm(100000)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}

		m := r.MultipartForm

		files := m.File["userfile"]
		for i, _ := range files {

			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			dst, err := os.Create("/upload/" + files[i].Filename)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}

		return
	} else {
		w.WriteHeader(500)
	}
}

func upload2server(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		// 上传页面
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(200)
		html := `
<form enctype="multipart/form-data" action="/upload" method="POST">
    Select file: <input name="userfile" type="file" maxlength="18" size="18" />
    <input type="submit" value="upload" />
</form>
`
		io.WriteString(w, html)
	} 
}

func deleteServer(w http.ResponseWriter, r *http.Request) {

	path := "/upload"

	now_time := time.Now().Unix()

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		file_time := f.ModTime().Unix()
		if (now_time - file_time) > diff_time {
			fmt.Printf("Delete file %v !\r\n", path)
			os.RemoveAll(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\r\n", err)
	}

	w.WriteHeader(200)
}

func main() {
	http.HandleFunc("/upload", uploadServer)
	http.HandleFunc("/upload2", upload2server)
	http.HandleFunc("/delete", deleteServer)
	http.Handle("/", http.FileServer(http.Dir("/")))
	err := http.ListenAndServe(":5360", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
