package util

import (
	"encoding/json"
	"fmt"
	"os"
)

type TargetNo struct {
	TargetNo string `json:"targetNo"`
}

func WriteJsonFile(targetNo string) {
	logFolderPath := "./data"
	logFilePath := fmt.Sprintf("%s/target_no.json", logFolderPath)
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		os.MkdirAll(logFolderPath,0777)
	}

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.Create(logFilePath)
	}
	data := make([]TargetNo, 1)
	data[0].TargetNo = targetNo

	println(data)
	doc, _ := json.Marshal(data)

	err := os.WriteFile(logFilePath, doc, os.FileMode(0644))
	if err != nil {
		WriteLogToFile(err.Error())
		return
	}
}

func ParseJsonFile() string {
	// JSON 파일 읽기
	file, err := os.ReadFile("./data/target_no.json")
	if err != nil {
		WriteLogToFile("파일을 읽는 도중 에러 발생:" + err.Error())
		panic(err)
	}

	// JSON 데이터 파싱하여 구조체에 저장
	var data []TargetNo
	err = json.Unmarshal(file, &data)
	if err != nil {
		WriteLogToFile("JSON 데이터 파싱 도중 에러 발생:" + err.Error())
		panic(err)
	}

	// 파싱된 데이터 사용 예시
	for _, item := range data {
		WriteLogToFile("TargetNo:" + item.TargetNo)
		return item.TargetNo
	}
	return ""
}