package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
)

var (
	logFileName = flag.String("log", "http_c.log", "Log file name")
	url         = flag.String("url", "http://192.168.1.200:7000/sms", "http url")
	bodyFile    = flag.String("body", "body.json", "http body data")
	qChan       chan int
	totalCount  int
)

func httpGet(url string, body string) (resp string, err error) {
	client := &http.Client{}
	bodySlice := []byte(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodySlice))
	if err != nil {
		return "", errors.New("new request failed")
	}

	req.Header.Add("Basic", "abcd")
	ret, err := client.Do(req)
	if err != nil {
		return "", errors.New("do request failed")
	}

	response, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return "", errors.New("get http response failed")
	}

	qChan <- 1
	return string(response), nil
}

func getHTTPBody(filename string) (httpBody string, err error) {
	dataFile, err := os.Open(filename)
	if err != nil {
		return "", errors.New("open http data file failed")
	}

	defer dataFile.Close()

	content, err := ioutil.ReadAll(dataFile)
	if err != nil {
		return "", errors.New("read http data failed")
	}

	return string(content), nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// 读取http_body内容
	body, err := getHTTPBody(*bodyFile)
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	log.Printf("%s", body)

	//set logfile
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "Http Client Test start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	qChan = make(chan int)

	totalCount = 1
	for i := 0; i < totalCount; i++ {
		go httpGet(*url, body)
	}

	i := 0
	for {
		<-qChan
		i++
		if i >= totalCount {
			break
		}
	}
}
