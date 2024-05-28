package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/util"
)

func AddTarget(targetNo, targetGroupNo string) *string {

	var result *string

	apiKeys := ncloud.Keys()
	client := vloadbalancer.NewAPIClient(vloadbalancer.NewConfiguration(apiKeys))

	var targetNoList []string
	targetNoList = append(targetNoList, targetNo)

	util.WriteLogToFile("타겟 추가 요청")
	addReq := vloadbalancer.AddTargetRequest{
		TargetGroupNo: ncloud.String(targetGroupNo),
		TargetNoList:  ncloud.StringList(targetNoList),
	}

	if r, err := client.V2Api.AddTarget(&addReq); err != nil {
		util.WriteLogToFile(err.Error())
		println("ERR")
	} else {
		util.WriteLogToFile(ncloud.StringValue(r.RequestId))
		util.WriteLogToFile(ncloud.StringValue(r.ReturnCode))
		util.WriteLogToFile(ncloud.StringValue(r.ReturnMessage))
		println("DONE")
		result = r.ReturnCode
	}
	return result
}
