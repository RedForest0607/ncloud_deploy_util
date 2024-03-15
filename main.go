package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"
	"time"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vloadbalancer"
)

func ParseTargetNo(targetName string, jsonData []byte) (string, error) {
    type Target struct {
        TargetNo   string `json:"targetNo"`
        TargetName string `json:"targetName"`
    }

    type GetTargetListResponse struct {
        TotalRows   int      `json:"totalRows"`
        TargetList  []Target `json:"targetList"`
        RequestId   string   `json:"requestId"`
        ReturnCode  string   `json:"returnCode"`
        ReturnMessage string `json:"returnMessage"`
    }

    var response struct {
        GetTargetListResponse GetTargetListResponse `json:"getTargetListResponse"`
    }

    err := json.Unmarshal(jsonData, &response)
    if err != nil {
        return "", err
    }

    for _, target := range response.GetTargetListResponse.TargetList {
        if target.TargetName == targetName {
            return target.TargetNo, nil
        }
    }

    return "", fmt.Errorf("target with name %s not found", targetName)
}
func MakeSignature(accessKey, secretKey, url string, timestamp int64) string {
    space := " "
    newLine := "\n"
    method := "GET"

    message := method + space + url + newLine + fmt.Sprintf("%d", timestamp) + newLine + accessKey

    key := []byte(secretKey)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))

    hash := h.Sum(nil)

    return base64.StdEncoding.EncodeToString(hash)
}

func CreateRequest(accessKey, signature string, timestamp int64) http.Header {
    headers := make(http.Header)

    headers.Set("x-ncp-apigw-timestamp", fmt.Sprintf("%d", timestamp))
    headers.Set("x-ncp-iam-access-key", accessKey)
    headers.Set("x-ncp-apigw-signature-v2", signature)

    return headers
}

func GetTargetNo(targetName,targetGroupNo string, apiKey *ncloud.APIKey) string {
    // 현재 시간을 밀리초로 변환하여 정수로 반올림
    // 정수를 문자열로 변환
    timestamp := time.Now().Unix() * 1000

    accessKey := apiKey.AccessKey
    secretKey := apiKey.SecretKey

    apigwUrl := "https://ncloud.apigw.ntruss.com"
    queryUrl := "/vloadbalancer/v2/getTargetList?regionCode=KR&targetGroupNo="+targetGroupNo
    responseType := "&responseFormatType=json"
    url := apigwUrl + queryUrl + responseType

    signature := MakeSignature(accessKey, secretKey, queryUrl + responseType, timestamp)
    headers := CreateRequest(accessKey, signature, timestamp)

    // HTTP 요청 생성
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "error"
    }

    // 생성한 헤더 설정
    req.Header = headers

    // HTTP 클라이언트 생성
    client := &http.Client{}

    // 요청 보내기
    resp, err := client.Do(req)
    if err != nil {
        return "error"
    }
    defer resp.Body.Close()

    // 응답 읽기
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "error"
    }

    // 응답 결과 변환
    targetNo, err := ParseTargetNo(targetName, body)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        return targetNo
    }

    return "not_found"
}

func main() {

    log.Println("\n######## REQUEST TARGET LIST ###########")
    targetGroupNo := os.Getenv("TARGET_GROUP_NO")
	targetName := os.Getenv("TARGET_NAME")
	apiKeys := ncloud.Keys()

	result := GetTargetNo(targetName, targetGroupNo, apiKeys)
	log.Println(result)
	client := vloadbalancer.NewAPIClient(vloadbalancer.NewConfiguration(apiKeys))

    var targetNoList []string
    targetNoList = append (targetNoList, result)

	log.Println("\n######## REMOVE TARGET ###########")
	removeReq := vloadbalancer.RemoveTargetRequest{
		TargetGroupNo: ncloud.String(targetGroupNo),
		TargetNoList: ncloud.StringList(targetNoList),
	}

	if r, err := client.V2Api.RemoveTarget(&removeReq); err != nil {
		log.Println(err)
	} else {
		log.Println(ncloud.StringValue(r.RequestId))
		log.Println(ncloud.StringValue(r.ReturnCode))
		log.Println(ncloud.StringValue(r.ReturnMessage))
	}

	time.Sleep(time.Second * 10)

	log.Println("\n######## ADD TARGET ###########")
	addReq := vloadbalancer.AddTargetRequest{
		TargetGroupNo: ncloud.String(targetGroupNo),
		TargetNoList: ncloud.StringList(targetNoList),
	}

	if r, err := client.V2Api.AddTarget(&addReq); err != nil {
		log.Println(err)
	} else {
		log.Println(ncloud.StringValue(r.RequestId))
		log.Println(ncloud.StringValue(r.ReturnCode))
		log.Println(ncloud.StringValue(r.ReturnMessage))
	}
}