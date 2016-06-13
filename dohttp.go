package main

import "net/http"

// func httpPost(url string, body string, idx int) {
// 	client := &http.Client{}
// 	bodySlice := []byte(body)

// 	specURL := url + "?name=" + strconv.Itoa(idx)
// 	req, err := http.NewRequest("POST", specURL, bytes.NewBuffer(bodySlice))
// 	if err != nil {
// 		log.Printf("new request failed")
// 		return
// 	}

// 	req.Header.Add("Basic", "abcd")
// 	ret, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("do request failed")
// 		return
// 	}

// 	response, err := ioutil.ReadAll(ret.Body)
// 	if err != nil {
// 		log.Printf("get http response failed")
// 		return
// 	}

// 	log.Printf("idx = %d, %s\n", idx, string(response))
// }

func httpGet(url string, statChan chan int) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger.Errorf("new request failed, err = %s", err)
		return
	}

	req.Header.Add("Basic", "abcd")
	response, err := client.Do(req)
	if err != nil {
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
			go httpGet(testCase.URL, statChan)
		}
	default:
		Logger.Errorf("Method is invalid, err = %s", testCase.Method)
	}
}
