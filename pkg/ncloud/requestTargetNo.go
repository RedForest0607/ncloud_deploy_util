package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/util"
)

func RequestInstanceNo(targetGroupNo string) {

	apiKeys := ncloud.Keys()
	client := vloadbalancer.NewAPIClient(vloadbalancer.NewConfiguration(apiKeys))

	var targetGroupNoList []string
	targetGroupNoList = append(targetGroupNoList, targetGroupNo)

	util.WriteLogToFile("\n######## REQUEST TARGET LIST ###########")
	addReq := vloadbalancer.GetTargetGroupListRequest{
		TargetGroupNoList: ncloud.StringList(targetGroupNoList),
	}

	if r, err := client.V2Api.GetTargetGroupList(&addReq); err != nil {
		util.WriteLogToFile(err.Error())
	} else {
		fmt.Println(r)
		data := r.TargetGroupList[0].TargetNoList
		for _, item := range data {
			fmt.Printf(*item + "\n")
		}
	}
}

func RequestTargetNo(instanceName string) string {
	apiKeys := ncloud.Keys()
	client := vserver.NewAPIClient(vserver.NewConfiguration(apiKeys))

	req := vserver.GetServerInstanceListRequest{
		ServerName: ncloud.String(instanceName),
	}

	if r, err := client.V2Api.GetServerInstanceList(&req); err != nil {
		util.WriteLogToFile(err.Error())
	} else {
		data := r.ServerInstanceList
		for _, item := range data {
			util.WriteLogToFile("타겟 번호 조회 완료")
			util.WriteLogToFile(*item.ServerName + "     " + *item.ServerInstanceNo + "\n")
			return *item.ServerInstanceNo
		}
	}
	return ""
}
