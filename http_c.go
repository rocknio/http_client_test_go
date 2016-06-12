package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

var (
	logFileName = flag.String("log", "http_c.log", "Log file name")
	url         = flag.String("url", "http://192.168.1.8:9000/", "http url")
	bodyFile    = flag.String("body", "body.json", "http body data")
	qChan       chan int
	totalCount  int
)

type smsContent struct {
	Template  string `json:"template"`
	DestTelno string `json:"dest_telno"`
	Content   string `json:"content"`
}

func httpPost(url string, body string, idx int) {
	client := &http.Client{}
	bodySlice := []byte(body)

	specURL := url + "?name=" + strconv.Itoa(idx)
	req, err := http.NewRequest("POST", specURL, bytes.NewBuffer(bodySlice))
	if err != nil {
		log.Printf("new request failed")
		qChan <- 1
		return
	}

	req.Header.Add("Basic", "abcd")
	ret, err := client.Do(req)
	if err != nil {
		log.Printf("do request failed")
		qChan <- 1
		return
	}

	response, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		log.Printf("get http response failed")
		qChan <- 1
		return
	}

	log.Printf("idx = %d, %s\n", idx, string(response))
	qChan <- 1
}

func httpGet(url string, idx int) {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("http get failed")
		return
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read http response failed")
		return
	}

	log.Printf("idx = %d, %s\n", idx, string(resBody))
	qChan <- 1
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

func parseJSON(content string) {
	var tmp smsContent
	json.Unmarshal([]byte(content), &tmp)
	fmt.Printf("%v", tmp)

	var f interface{}
	json.Unmarshal([]byte(content), &f)
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Printf("%v is %v", k, vv)

		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
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

	//set logfile
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "Http Client Test start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	log.Printf("%s\n", body)
	parseJSON(body)

	qChan = make(chan int)

	totalCount = 100
	for i := 0; i < totalCount; i++ {
		// go httpGet(*url, i)
		go httpPost(*url, body, i)
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
