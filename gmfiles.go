/* 2021/06/11 */

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var FileList []string
var Path string

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

func getter() {

	listernner, err := net.Listen("tcp", "127.0.0.1:5361")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	defer listernner.Close()

	fmt.Println("sending for ....")
	conn, err := listernner.Accept()
	if err != nil {
		fmt.Println("listenner.Accept err=", err)
		return
	}

	defer conn.Close()

	for i := 0; i < len(FileList); i++ {

		sender(FileList[i], conn)

		if i == len(FileList)-1 {
			conn.Write([]byte("over"))
			fmt.Println("The task is over !")
			return
		}
	}
}

func sender(TmFileList string, conn net.Conn) {

	info, err := os.Stat(TmFileList)
	if err != nil {
		WriteLog("os.Stat err= ", err.Error())
		return
	}

	_, err = conn.Write([]byte(info.Name()))
	if err != nil {
		WriteLog("conn.Write err =", err.Error())
		return
	}
	fmt.Println("send ", TmFileList)
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
		fmt.Println("os.Open err=", err)
		WriteLog("os.Open err=", err.Error())
		return
	}
	defer f.Close()

	buf := make([]byte, 1024*16)
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

		if n == 0 {
			fmt.Println("n==0 File received ! ")
			return
		}

		conn.Write(buf[:n])
	}
}

func main() {

	FileList, _ = GetFileList(".", FileList)
	listernner, err := net.Listen("tcp", "127.0.0.1:5360")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	defer listernner.Close()

	for {
		fmt.Println("Waiting for ....")
		conn, err := listernner.Accept()
		if err != nil {
			fmt.Println("listenner.Accept err1 =", err)
		}
		defer conn.Close()

		var n = 0
		buf := make([]byte, 1024)
		n, err = conn.Read(buf)
		if err != nil {
			WriteLog("conn.Read err=", err.Error())
			return
		}

		commands := string(buf[:n])
		switch {
		case commands == "ls":
			go Listfiles(conn)
		case commands == "put":
			go putter()
		case commands == "get":
			go getter()
		}

	}
}

func Listfiles(conn net.Conn) {

	_, err1 := conn.Write([]byte(strings.Join([]string(FileList), ",")))
	if err1 != nil {
		WriteLog("listfiles conn.Write err1 =", err1.Error())
		return
	}
	buf := make([]byte, 1024)
	n, err2 := conn.Read(buf)
	if err2 != nil {
		WriteLog("listfiles conn.Read err2=", err2.Error())
		return
	}
	if "over" == string(buf[:n]) {
		fmt.Println("select  complete !")
	}
	return
}

func isDirExists(pathname string) bool {
	_, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir() (pathTmp string) {

	Path, _ = os.Getwd()
	//t := time.Now()
	//date := t.Format("20060102")
	pathTmp = Path + "/ "       
	// + date + "/" 
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

func putter() {

	pathTmp := CreateDir()
	listernner, err := net.Listen("tcp4", "127.0.0.1:5362")
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	defer listernner.Close()

	fmt.Println("putting for ....")
	conn, err := listernner.Accept()
	if err != nil {
		fmt.Println("listenner.Accept err=", err)
	}
	defer conn.Close()

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
		RecvFile(pathTmp, filename, conn)
	}
}

func RecvFile(pathTmp string, filename string, conn net.Conn) {

	f, err := os.OpenFile(pathTmp+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
			fmt.Println("n==0 File received ")
			return
		}

		if "complete" == string(buf[:n]) {
			fmt.Println("complete")
			return
		}

		f.Write(buf[:n])
	}

}
