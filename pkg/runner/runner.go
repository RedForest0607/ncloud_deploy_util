package runner

import (
	"encoding/json"
	"fmt"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/ncloud"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/util"
	"log"
	"os"
)

func Start() {
	targetGroupNo := os.Getenv("TARGET_GROUP_NO")
	targetName := os.Getenv("TARGET_NAME")
	behavior := os.Getenv("BEHAVIOR")

	if behavior == "ADD" {
		targetNo := util.ParseJsonFile(targetName)
		if targetNo == "" {
			util.WriteLogToFile("NO TARGET NO")
			println("ERR")
			return
		}
		ncloud.AddTarget(targetNo, targetGroupNo)
	} else if behavior == "REMOVE" {
		targetNo := ncloud.SendTargetNoRequest(targetName, targetGroupNo)
		if targetNo == "" || targetNo == "not_found" {
			util.WriteLogToFile("TARGET NO NOT FOUND")
			println("ERR")
			return
		}
		util.WriteJsonFile(targetName, targetNo)
		ncloud.RemoveTarget(targetNo, targetGroupNo)
	} else if behavior == "FIND" {
		filePath := "./data/target_no.json"
		ncloud.RequestInstanceNo(targetName)
		// 파일 존재 여부 확인
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			util.WriteLogToFile("target_no.json 파일을 생성합니다.")
			_, err := os.Create("./data/target_no.json")

			// JSON 객체 생성
			emptyJSON := make(map[string]interface{})
			// JSON 데이터 마샬링
			jsonBytes, err := json.MarshalIndent(emptyJSON, "", "  ")
			if err != nil {
				log.Fatalf("JSON 데이터를 마샬링하는 중 오류가 발생했습니다: %v", err)
			}

			// 파일에 JSON 쓰기
			if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
				log.Fatalf("파일에 JSON을 쓰는 중 오류가 발생했습니다: %v", err)
			}
		} else if err != nil {
			fmt.Printf("파일을 확인하는 중 오류가 발생했습니다: %v\n", err)
			return
		}
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		var jsonData map[string]interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			log.Fatalf("JSON 데이터를 파싱하는 중 오류가 발생했습니다: %v", err)
		}
		value, exists := jsonData[targetName]
		if exists {
			log.Println(value)
		} else {
			targetNo := ncloud.SendTargetNoRequest(targetName, targetGroupNo)
			util.WriteJsonFile(targetName, targetNo)
		}

	} else {
		targetNo := ncloud.SendTargetNoRequest(targetName, targetGroupNo)
		ncloud.RequestTargetNo(targetGroupNo)
		util.WriteJsonFile(targetName, targetNo)
	}
}
