package ncloud

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/util"
)

func RemoveTarget(targetNo, targetGroupNo string) *string {

	var result *string
	apiKeys := ncloud.Keys()

	client := vloadbalancer.NewAPIClient(vloadbalancer.NewConfiguration(apiKeys))

	var targetNoList []string
	targetNoList = append(targetNoList, targetNo)

	util.WriteLogToFile("타겟 삭제 요청")
	removeReq := vloadbalancer.RemoveTargetRequest{
		TargetGroupNo: ncloud.String(targetGroupNo),
		TargetNoList:  ncloud.StringList(targetNoList),
	}

	if r, err := client.V2Api.RemoveTarget(&removeReq); err != nil {
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
