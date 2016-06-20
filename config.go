package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type testInfo struct {
	StartTime time.Time // 测试开始时间
	EndTime   time.Time // 测试结束时间
	SendCount int       // 已发送请求数
	SuccCount int       // 成功次数
	FailCount int       // 失败次数
}

// TestCase 单个测试用例结构
type TestCase struct {
	TestName       string   `json:"test_name"`        // 测试名
	URL            string   `json:"url"`              // 测试URL
	Basic          string   `json:"basic"`            // http头basic认证字符串
	Method         string   `json:"method"`           // http方法
	PostContent    string   `json:"post_content"`     // 如果是post方法，post的内容
	TotalCount     int      `json:"total_count"`      // 总发送请求参数
	CountPerSecond int      `json:"count_per_second"` // 每秒发送请求次数
	TestInfo       testInfo // 测试过程信息
}

func readConfig(fileName string) (content string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		Logger.Errorf("Open the configfile filed, err = %s", err)
		return "", err
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		Logger.Errorf("Open the configfile filed, err = %s", err)
		return "", err
	}

	content = string(fileContent)
	err = nil
	return
}

// ParseConfig 解析JSON格式配置文件，返回结构体slice
func ParseConfig(fileName string) (testCase TestCase, err error) {
	config, err := readConfig(fileName)
	err = json.Unmarshal([]byte(config), &testCase)
	if err != nil {
		Logger.Errorf("parse config failed! err = %s", err)
	}

	return
}
