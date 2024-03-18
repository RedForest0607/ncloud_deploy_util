package ncloud

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"lab.overpass.co.kr/aws/ncloud-deployer/pkg/util"
)

const ApigwUrl string = "https://ncloud.apigw.ntruss.com"
const QueryUrl string = "/vloadbalancer/v2/getTargetList?regionCode=KR&targetGroupNo="

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

    targetListQueryUrl := QueryUrl + targetGroupNo
    responseType := "&responseFormatType=json"
    url := ApigwUrl + QueryUrl + targetGroupNo + responseType

    signature := MakeSignature(accessKey, secretKey, targetListQueryUrl + responseType, timestamp)
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
        util.WriteLogToFile(err.Error())
    } else {
        return targetNo
    }
    util.WriteLogToFile(string(body))
    return "not_found"
}

func SendTargetNoRequest(targetName, targetGroupNo string) string {

    util.WriteLogToFile("\n######## REQUEST TARGET LIST ###########")
	apiKeys := ncloud.Keys()
	result := GetTargetNo(targetName, targetGroupNo, apiKeys)
	util.WriteLogToFile(result)

    return result
}