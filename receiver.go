package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func isDirExists(pathname string) bool {

	_, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func RecvFile(pathTmp string, filename string, conn net.Conn) {

	f, err := os.OpenFile(pathTmp+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("os.Create err=", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)
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
			break
		}
		f.Write(buf[:n])
	}
}

func ReceiveHandle(pathTmp string, conn net.Conn) {

	defer conn.Close()
	buf := make([]byte, 1024*16)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err =", err)
		return
	}
	filename := string(buf[:n])
	conn.Write([]byte("ok"))
	fmt.Println("Prereceive file :", filename)
	RecvFile(pathTmp, filename, conn)
}

func CreateDir() (pathTmp string) {

	Path, _ := os.Getwd()
	t := time.Now()
	date := t.Format("20060102")
	pathTmp = Path + "/" + date + "/"
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

func main() {

	pathTmp := CreateDir()
	listernner, err := net.Listen("tcp4", "0.0.0.0:5360")
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
			continue
		}
		go ReceiveHandle(pathTmp, conn)
	}
}
