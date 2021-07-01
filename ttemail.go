package main

import (
    "fmt"
    "log"
    "net/smtp"
    "goemail"
	"io/ioutil"
	"time"
	"strings"
	"os"
)

var FileList []string
var pathTmp string

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

func isDirExists(pathname string) bool {
	_, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func searchDir()(pathTmp string){

	Path, _ := os.Getwd()
	t := time.Now()
	date := t.Format("20210701")
	pathTmp = Path + "/" + date + "/" 
	if isDirExists(pathTmp) {
		fmt.Println("Directory exists ")
	} else {
		fmt.Println("Directory does not exist")		
	}
	return pathTmp
}

func main() {

    pathTmp = searchDir()

    FileList, _ = GetFileList(pathTmp, FileList)
	
	sendemail(FileList) 
}

func sendemail (FileList []string ){

    e := email.NewEmail()  
	 
    e.From = "gecx1057@163.com"

    e.To = []string{"chunnet@139.com"}
 
    e.Subject = "ceshi"

    e.Text = []byte(strings.Join([]string(FileList), ","))
	
    err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "gecx1057", "UAQMGIOJLTJJHDCL", "smtp.163.com"))
    if err != nil {
        log.Fatal(err)
    }
    return
}


