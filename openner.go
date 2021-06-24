/* 2021/06/11 */

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var FileSize float64
var Addrstring string
var putAddrstring string
var getAddrstring string
var FileList []string
var commadns string

func putter() {

	FileList, _ = GetFileList(".", FileList)

	conn, err := net.Dial("tcp4", putAddrstring)
	if err != nil {
		WriteLog("net.Dial err=", err.Error())
		return
	}
	defer conn.Close()

	for i := 0; i < len(FileList); i++ {

		putfile(FileList[i], conn)

		if i == len(FileList)-1 {
			conn.Write([]byte("over"))
			fmt.Println("The task is over !")
			return
		}
	}
}

func putfile(TmFileList string, conn net.Conn) {

	info, err := os.Stat(TmFileList)
	if err != nil {
		WriteLog("os.Stat err= ", err.Error())
		return
	}

	fmt.Printf("Send FileName : %s ", info.Name())
	fmt.Println("FileSize : ", info.Size())
	FileSize64 := info.Size()
	strInt64 := strconv.FormatInt(FileSize64, 10)
	FileSize, err = strconv.ParseFloat(strInt64, 64)

	_, err = conn.Write([]byte(info.Name()))
	if err != nil {
		WriteLog("conn.Write err =", err.Error())
		return
	}

	var n = 0
	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		WriteLog("conn.Read err=", err.Error())
		return
	}

	if "ok" == string(buf[:n]) {
		SendFile(TmFileList, conn)
	}

	conn.Write([]byte("complete"))
	return
}

func SendFile(filename string, conn net.Conn) {

	f, err := os.Open(filename)
	if err != nil {
		WriteLog("os.Open err=", err.Error())
		return
	}
	defer f.Close()

	var sum float64
	sum = 0
	buf := make([]byte, 1024*16)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File transfer complete:", filename)
				WriteLog("File transfer complete:", filename)
			} else {
				WriteLog(" f.Read err=", err.Error())
			}
			return

			if n == 0 {
				fmt.Println("n==0 File complete! ")
				return
			}
		}

		conn.Write(buf[:n])
		sum += float64(n)
		fmt.Printf("\r%.0f%%", sum/FileSize*100)
	}
}

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
		WriteLog("read dir fail:", err.Error())
		return FileList, err
	}

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		FileList = append(FileList, fi.Name())
	}
	return FileList, nil
}

func isDirExists(pathname string) bool {
	_, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir() (pathTmp string) {

	Path, _ := os.Getwd()
	t := time.Now()
	date := t.Format("20060102")
	pathTmp = Path + "/ " + date + "/"
	if isDirExists(pathTmp) {
		fmt.Println("Directory exists ")
	} else {
		//fmt.Println("Directory does not exist
		err := os.Mkdir(pathTmp, 0777)
		if err != nil {
			fmt.Println("CreateDir err", err)
		}
	}
	return pathTmp
}

func init() {

	var host string
	var port int
	flag.StringVar(&host, "h", "localhost", "")
	flag.IntVar(&port, "port", 5360, "")
	flag.Parse()
	Addrstring = host + ":" + strconv.Itoa(port)
	getAddrstring = host + ":" + "5361"
	putAddrstring = host + ":" + "5362"

}

func getter() {

	pathTmp := CreateDir()

	conn, err := net.Dial("tcp", getAddrstring)
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read Gerr =", err)
			return
		}

		if "over" == string(buf[:n]) {
			fmt.Println("The task is over !")
			return
		}

		filename := string(buf[:n])
		fmt.Println("Prereceive file :", filename)
		conn.Write([]byte("ok"))
		DownFile(conn, pathTmp, filename)
	}
}

func DownFile(conn net.Conn, DpathTmp string, Dfilename string) {

	f, err := os.OpenFile(DpathTmp+"/"+Dfilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("os.Create err=", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024*16)
	for {

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File received ")
			} else {
				fmt.Println("conn.Read err=", err)
			}
			return
		}

		if n == 0 {
			fmt.Println("n==0 File received ! ")
			return
		}
		if "complete" == string(buf[:n]) {
			fmt.Println("complete")
			return
		}

		f.Write(buf[:n])
	}
}

func Listfiles(conn net.Conn) {

	buf := make([]byte, 1024)
	n, err1 := conn.Read(buf)
	if err1 != nil {
		fmt.Println("read  fail:", err1)
	}
	_, err3 := conn.Write([]byte("over"))
	if err3 != nil {
		fmt.Println("Listfiles conn.Write err3", err3)
		return
	}
	FileList := string(buf[:n])

	fmt.Println(FileList)
	return
}

func main() {

	for {
		conn, err := net.Dial("tcp", Addrstring)
		if err != nil {
			fmt.Println("net.Dial err=", err)
		}
		//defer conn.Close()

		fmt.Println("put | get | ls | exit ")
		fmt.Scan(&commadns)
		_ = make([]byte, 1024)
		_, err = conn.Write([]byte(commadns))
		if err != nil {
			WriteLog("conn.Write err =", err.Error())
			return
		}

		switch {
		case commadns == "ls":
			Listfiles(conn)
		case commadns == "put":
			putter()
			return
		case commadns == "get":
			getter()
			return
		case commadns == "exit":
			return
		}
	}
}
