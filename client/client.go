package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	addr        = "localhost"
	port        = ":5000"
	getRoute    = "/file"
	postRoute   = "/upload"
	queryParams = "?filename="
	filename    = os.Getenv("FILENAME")
)

func main() {
	route := "http://" + addr + port
	//PostFile(route)
	GetFile(route)
}

func PostFile(route string) {
	req := route + postRoute + queryParams + filename
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, _ := bodyWriter.CreateFormFile("filename", filename)

	// open file handle
	fh, err := os.Open("./" + filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fh.Close()

	//iocopy
	io.Copy(fileWriter, fh)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(req, contentType, bodyBuf)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	bd, _ := ioutil.ReadAll(resp.Body)
	os.Stdout.Write(bd)
	fmt.Println("Succesfully stored file")
}

func GetFile(route string) {
	req := route + getRoute + queryParams + filename
	resp, err := http.Get(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	// if response is type json (error encountered), show in terminal
	if resp.Header.Get("Content-Type") != "application/octet-stream" {
		bd, _ := ioutil.ReadAll(resp.Body)
		os.Stdout.Write(bd)
		os.Stdout.Write([]byte("\n"))
		return
	}
	// Create blank file
	newFilePath := "./" + filename
	file, err := os.Create(newFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully downloaded file")
}
