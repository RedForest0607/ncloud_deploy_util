package runner

import (
	"os"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/ncloud"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/util"
)

func Start() {
	targetGroupNo := os.Getenv("TARGET_GROUP_NO")
    targetName := os.Getenv("TARGET_NAME")
    behavior := os.Getenv("BEHAVIOR")

    if behavior == "ADD" {
		targetNo := util.ParseJsonFile()
		if targetNo == "" {
			util.WriteLogToFile("NO TARGET NO")
			return
		}
		ncloud.AddTarget(targetNo, targetGroupNo)
    } else if behavior == "REMOVE" {
		targetNo := ncloud.SendTargetNoRequest(targetName, targetGroupNo)
		util.WriteJsonFile(targetNo)
		ncloud.RemoveTarget(targetNo, targetGroupNo)
	}
}