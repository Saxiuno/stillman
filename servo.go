package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

var FileList []string

func WriteLog(Tips string, error string) {

	file := "ErrorMessage" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "\r\n", log.Ldate|log.Ltime)
	logger.Print(Tips, error)
}

func GetFileList(pathname string, FileList []string) ([]string, error) {

	files, err := ioutil.ReadDir(pathname)

	if err != nil {
		fmt.Println("read dir fail:", err)
		return FileList, err
	}

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		println(fi.Name())
		FileList = append(FileList, fi.Name())
	}
	return FileList, nil
}

func main() {

	FileList, _ = GetFileList("./servo", FileList)

	listernner, err := net.Listen("tcp", "0.0.0.0:5360")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	defer listernner.Close()
	for {
		fmt.Println("Waiting for ....")
		conn, err := listernner.Accept()
		if err != nil {
			fmt.Println("listenner.Accept err=", err)
		}
		defer conn.Close()

		var wg sync.WaitGroup
		wg.Add(1)
		go servo(&wg, conn)
		wg.Wait()
	}
}

func servo(wg *sync.WaitGroup, conn net.Conn) {

	for _, Filestr := range FileList {

		var n = 0
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			WriteLog("conn.Read err=", err.Error())
		}

		if "complete" == string(buf[:n]) {

			info, err := os.Stat(Filestr)
			if err != nil {
				WriteLog("os.Stat err= ", err.Error())
				return
			}

			_, err = conn.Write([]byte(info.Name()))
			if err != nil {
				WriteLog("conn.Write err =", err.Error())
				return
			}
			SendFile(Filestr, conn)
		}
	}
	conn.Write([]byte("!chunnet@qq.com!"))
	wg.Done()
}

func SendFile(filename string, conn net.Conn) {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("os.Open err=", err)
		WriteLog("os.Open err=", err.Error())
		return
	}
	defer f.Close()

	buf := make([]byte, 1024*4)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File transfer complete")
			} else {
				WriteLog(" f.Read err=", err.Error())
			}
			return
		}
		conn.Write(buf[:n])
	}
}
