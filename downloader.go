package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

var Addrstring string

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
		err := os.Mkdir(pathTmp, 0777)
		if err != nil {
			fmt.Println("CreateDir err", err)
		}
	}
	return pathTmp
}

func init() {

	var user string
	var pwd string
	var host string
	var port int
	flag.StringVar(&user, "u", "", "用户名,默认为空")
	flag.StringVar(&pwd, "pwd", "", "密码,默认为空")
	flag.StringVar(&host, "h", "localhost", "")
	flag.IntVar(&port, "port", 5360, "")
	flag.Parse()
	Addrstring = host + ":" + strconv.Itoa(port)
}

func main() {

	pathTmp := CreateDir()

	conn, err := net.Dial("tcp", Addrstring)
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()

	for {
		conn.Write([]byte("complete"))
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err =", err)
			return
		}

		filename := string(buf[:n])
		fmt.Println("Prereceive file :", filename)
		download(pathTmp, filename, conn)

	}
}

func download(pathTmp string, filename string, conn net.Conn) {

	f, err := os.OpenFile(pathTmp+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("os.Create err=", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
        f.Write(buf[:n])
		if n == 0 {
			fmt.Println("n==0 File received ")			
			break
		} 

        if n < 1024 {
			fmt.Println("n<1024 File received2 ")		
			break
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println("File received ")
				break
			} else {
				fmt.Println("conn.Read err=", err)
				break
			}
		}
	}
	return
}
