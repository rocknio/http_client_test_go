package main

import (
	"bytes"
	"net/http"
)

func httpPost(url string, body string, basicStr string, statChan chan int) {
	client := &http.Client{}
	bodySlice := []byte(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodySlice))
	if err != nil {
		statChan <- 500
		Logger.Errorf("new request failed, err = %s", err)
		return
	}

	authStr := "Basic " + basicStr
	req.Header.Add("Authorization", authStr)
	Logger.Infof("*****send post*****")
	response, err := client.Do(req)
	if err != nil {
		Logger.Errorf("do request failed, err = %s", err)
		statChan <- 500
		return
	}

	Logger.Infof("*****recv post response*****")

	// 关闭连接
	defer response.Body.Close()

	// ret, _ := ioutil.ReadAll(response.Body)
	// fmt.Printf("url = %s, ret = %s\n", url, string(ret))
	statChan <- response.StatusCode
}

func httpGet(url string, basicStr string, statChan chan int) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		statChan <- 500
		Logger.Errorf("new request failed, err = %s", err)
		return
	}

	authStr := "Basic " + basicStr
	req.Header.Add("Authorization", authStr)
	response, err := client.Do(req)
	if err != nil {
		statChan <- 500
		Logger.Errorf("do request failed, err = %s", err)
		return
	}

	// 关闭连接
	defer response.Body.Close()

	statChan <- response.StatusCode
}

// DoTest 执行http测试用例
func DoTest(testCase TestCase, statChan chan int) {
	switch testCase.Method {
	case "GET", "get", "Get":
		for i := 0; i < (testCase.CountPerSecond / 10); i++ {
			go httpGet(testCase.URL, testCase.Basic, statChan)
		}
	case "POST", "post", "Post":
		for i := 0; i < (testCase.CountPerSecond / 10); i++ {
			URL := testCase.URL // + "?name=" + strconv.Itoa(i)
			go httpPost(URL, testCase.PostContent, testCase.Basic, statChan)
		}
	default:
		Logger.Errorf("Method is invalid, err = %s", testCase.Method)
	}
}
