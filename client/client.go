package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var (
	addr        = "localhost"
	port        = ":5000"
	getRoute    = "/file"
	postRoute   = "/upload"
	queryParams = "?filename="
	token       = os.Getenv("token")
	filename    = os.Getenv("FILENAME")
)

func main() {
	route := "http://" + addr + port
	//PostFile(route)
	GetFile(route)
}

func PostFile(route string) {
	url := route + postRoute + queryParams + filename
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

	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client := http.Client{Timeout: time.Second * 10}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	bd, _ := ioutil.ReadAll(resp.Body)
	os.Stdout.Write(bd)
	fmt.Println("Succesfully stored file")
}

func GetFile(route string) {
	url := route + getRoute + queryParams + filename
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
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
