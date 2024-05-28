package runner

import (
	"fmt"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/ncloud"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/util"
	"os"
)

func Start() {
	targetName := os.Getenv("TARGET_NAME")
	behavior := os.Getenv("BEHAVIOR")

	if behavior == "ADD" {
		util.WriteLogToFile("-------------타겟 추가 시작")
		targetGroupNo := os.Getenv("TARGET_GROUP_NO")
		targetNo := ncloud.RequestTargetNo(targetName)
		if targetNo == "" || targetNo == "not_found" {
			util.WriteLogToFile("TARGET NO NOT FOUND")
			println("ERR")
			return
		}
		ncloud.AddTarget(targetNo, targetGroupNo)
	} else if behavior == "REMOVE" {
		util.WriteLogToFile("-------------타겟 삭제 시작")
		targetGroupNo := os.Getenv("TARGET_GROUP_NO")
		targetNo := ncloud.RequestTargetNo(targetName)
		if targetNo == "" || targetNo == "not_found" {
			util.WriteLogToFile("TARGET NO NOT FOUND")
			println("ERR")
			return
		}
		ncloud.RemoveTarget(targetNo, targetGroupNo)
	} else if behavior == "FIND" {
		targetNo := ncloud.RequestTargetNo(targetName)
		fmt.Println(targetNo)
	} else {
		fmt.Println("동작 상태 \"BEHAVIOR\" 값을 지정해주세요")
	}
}
