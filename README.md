# NCLOUD에서 STRANGER 배포 시 사용할 수 있는 유틸 입니다.

### 환경설정  
`~/.ncloud/configure`  
위의 경로에 해당하는 `configure`설정파일을 생성하고 아래의 내용을 작성해 주어야 합니다.  
```
ncloud_access_key_id={네이버 클라우드에서 받은 값}
ncloud_secret_access_ke={네이버 클라우드에서 받은 값}
```  
네이버 클라우드 측에서는 대문자로 샘플을 게시해뒀지만 소문자로 해야 정상적으로 동작합니다  
  
GO를 통해서 실행시키는 샘플입니다  
`TARGET_GROUP_NO="1234567" TARGET_NAME="prod-hamonica-node" BEHAVIOR="ADD" go run ./main.go`  
바이너리 파일을 통해서 실행시키는 샘플입니다  
`TARGET_GROUP_NO="1234567" TARGET_NAME="prod-hamonica-node" BEHAVIOR="ADD" ./main`
### 빌드
`go build main.go`  

### 환경변수  
`TARGET_GROUP_NO` : 로드밸런서에서 삭제할 타겟이 존재하는 그룹의 번호 입니다  
`TARGET_NAME` : 삭제할 인스턴스의 이름입니다  
`BEHAVIOR` : 동작을 정의합니다  
- `REMOVE` : 인스턴스를 타겟그룹에서 삭제합니다  
- `ADD` : 인스턴스를 타겟그룹에 추가합니다  
  
  
### 주의점  
네이버 API를 통해서는 일반적인 방법으로 인스턴스의 타겟번호를 알아내기 어렵습니다. 그래서 이 유틸리티도 `삭제 후 추가`로 진행한다는 가정하에 작동합니다. data 폴더에 생성되는 json파일에서 타겟 번호를 알아내기 때문에, 해당 파일 내용 없이 `REMOVE`동작이 동작하지 않습니다.
  
추후 REMOVE에서도 타깃을 검색할 수 있도록 수정할 예정 입니다.
