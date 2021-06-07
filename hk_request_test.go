package hk_artemis_sdk

import (
	"fmt"
	"testing"
	"time"
)

const host = "https://127.0.0.1"
const path = "/artemis/api/artemis/v1/agreementService/securityParam/appKey/"
const (
	appKey    = "1234567"
	appSecret = "Qqg9BfQWFszAaQFKxsGh"
)

func Test_httpGet(t *testing.T) {
	aReq := ArtemisReq{
		Host:      host,
		Path:      path + appKey,
		AppKey:    appKey,
		AppSecret: appSecret,
	}
	got, err := aReq.HttpGet(nil, nil, nil, time.Second*10)
	if err != nil {
		panic(err)
	}
	fmt.Println(got)
}
