package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
	"flag"
	"strconv"
	"bytes"
)

var FileSize float64
var Addrstring string
var FileList []string

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

func CmdLine(Orders1 string, Orders2 string, Orders3 string) {

	command := exec.Command(Orders1, Orders2, Orders3)
	outinfo := bytes.Buffer{}
	command.Stdout = &outinfo
	err := command.Start()
	if err != nil {
		WriteLog("command.start()", err.Error())
	}
	if err = command.Wait(); err != nil {
		WriteLog("command.wait()", err.Error())
	} else{
	WriteLog("CmdLine Successsful",outinfo.String())
	}
	
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

func init(){

    var user string
    var pwd string
    var host string
    var port int
	var Addrstring string
    flag.StringVar(&user, "u", "", "用户名,默认为空")
    flag.StringVar(&pwd, "pwd", "", "密码,默认为空")
    flag.StringVar(&host, "h", "localhost","")
    flag.IntVar(&port, "port", 5360,"")
    flag.Parse()	
	Addrstring = host + ":"+ strconv.Itoa(port)
}

func stopCmdLine(){
    CmdLine("sc", "stop", "SQLSERVERAGENT")
	CmdLine("sc", "stop", "MSSQLSERVER")
	CmdLine("sc", "stop", "ReportServer")
	CmdLine("sc", "stop", "MSSQLFDLauncher")
	CmdLine("sc", "stop", "MSSQLServerOLAPService")
	CmdLine("sc", "stop", "MsDtsServer100")
	CmdLine("sc", "stop", "SQLWriter")
}

func startCmdLine(){
	CmdLine("sc", "start", "MSSQLSERVER")
	CmdLine("sc", "start", "SQLSERVERAGENT")
	CmdLine("sc", "start", "ReportServer")
	CmdLine("sc", "start", "MSSQLFDLauncher")
	CmdLine("sc", "start", "MSSQLServerOLAPService")
	CmdLine("sc", "start", "MsDtsServer100")
	CmdLine("sc", "start", "SQLWriter")
}

func main() {

	stoptCmdLine()
	
	FileList, _ = GetFileList(".", FileList)
	
	for i := 0; i < len(FileList); i++ {
		info, err := os.Stat(FileList[i])
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
			SendFile(FileList[i], conn)
			conn.Close()
		}
	}
    startCmdLine()	
}
