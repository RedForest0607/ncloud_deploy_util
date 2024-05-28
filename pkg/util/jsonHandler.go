package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJsonFile(targetName, targetNo string) {

	var jsonData map[string]interface{}
	logFolderPath := "./data"
	logFilePath := fmt.Sprintf("%s/target_no.json", logFolderPath)
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		os.MkdirAll(logFolderPath, 0777)
	}

	// 파일 존재 여부 확인
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		// 파일이 존재하지 않으면 새로운 맵 생성
		jsonData = make(map[string]interface{})
	} else {
		// 파일이 존재하면 파일 읽기
		data, err := os.ReadFile(logFilePath)
		if err != nil {
			WriteLogToFile(err.Error())
			return
		}

		// JSON 데이터 언마샬링
		if err := json.Unmarshal(data, &jsonData); err != nil {
			WriteLogToFile(err.Error())
			return
		}
	}

	jsonData[targetName] = targetNo

	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		WriteLogToFile(err.Error())
		return
	}

	err = os.WriteFile(logFilePath, jsonBytes, os.FileMode(0644))
	if err != nil {
		WriteLogToFile(err.Error())
		return
	}
}

func ParseJsonFile(targetName string) string {
	// JSON 파일 읽기
	file, err := os.ReadFile("./data/target_no.json")
	if err != nil {
		WriteLogToFile("파일을 읽는 도중 에러 발생:" + err.Error())
		panic(err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(file, &jsonData); err != nil {
		WriteLogToFile("파일을 읽는 도중 에러 발생:" + err.Error())
		panic(err)
	}

	// 파싱된 데이터 사용 예시
	value, exists := jsonData[targetName]
	if exists {
		return value.(string)
	} else {
		return ""
	}
}
