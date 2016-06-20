package main

import (
	"flag"
	"runtime"
	"time"
)

var (
	logFileName = flag.String("log", "test.log", "Log file name")
	cfgFileNmae = flag.String("config", "config.json", "Config file name")
)

func printTestInfo(testCase TestCase) {
	Logger.Infof("TotalCount = %d, Succ = %d, Fail = %d\n",
		testCase.TotalCount, testCase.TestInfo.SuccCount, testCase.TestInfo.FailCount)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// 日志设置
	err := InitLogger(*logFileName)
	if err != nil {
		return
	}

	// 读取配置文件
	testCase, err := ParseConfig(*cfgFileNmae)
	if err != nil {
		return
	}

	// 打印测试开始
	Logger.Info("*******************************************************************************************************\n")
	Logger.Infof("TEST = %s, Start......", testCase.TestName)

	statChan := make(chan int, 10000)

	// 启动一个协程，执行测试用例
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		testCase.TestInfo.StartTime = time.Now()
		for _ = range ticker.C {
			if testCase.TestInfo.SendCount < testCase.TotalCount {
				DoTest(testCase, statChan)
				testCase.TestInfo.SendCount += (testCase.CountPerSecond / 10)
			}
		}
	}()

	// 启动一个协程，每秒打印测试过程
	printTicker := time.NewTicker(time.Second)
	go func() {
		for _ = range printTicker.C {
			go printTestInfo(testCase)
		}
	}()

	//
	for {
		httpRespCode := <-statChan
		if httpRespCode == 200 {
			testCase.TestInfo.SuccCount++
		} else {
			testCase.TestInfo.FailCount++
		}

		if testCase.TotalCount == (testCase.TestInfo.SuccCount + testCase.TestInfo.FailCount) {
			break
		}
	}
	testCase.TestInfo.EndTime = time.Now()

	// 等待1秒，保证最后一次过程打印完成
	time.Sleep(1 * time.Second)

	Logger.Infof("\nTEST = %s\nTestInfo = %s,%s\nStartTime = %s\nEndTime   = %s\nTotalCount = %d, Succ = %d, Fail = %d\n",
		testCase.TestName, testCase.Method, testCase.URL, testCase.TestInfo.StartTime, testCase.TestInfo.EndTime,
		testCase.TotalCount, testCase.TestInfo.SuccCount, testCase.TestInfo.FailCount)
	Logger.Info("*******************************************************************************************************\n")
}
