package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func RecvFile(fullname string, conn net.Conn) {

    temppath :=filepath.Dir(fullname) 
    err2 :=os.MkdirAll(temppath,0766)
	if err2 != nil {
		fmt.Println("os.MkdirAll err=", err2)
		return
	}
	
	f, err := os.OpenFile(fullname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("os.Create err=", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024*64)
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

func ReceiveHandle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024*16)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err =", err)
		return
	}
	fullname := string(buf[:n])
	conn.Write([]byte("ok"))
	fmt.Println("Prereceive filename :", fullname)  
	RecvFile(fullname, conn)
}

func main() {

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
		go ReceiveHandle(conn)
	}
}
