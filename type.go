package hk_artemis_sdk

type ArtemisResp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
type Data struct {
	AppSecret  string `json:"appSecret"`
	Time       string `json:"time"`
	TimeSecret string `json:"timeSecret"`
}

type ArtemisReq struct {
	Host      string
	Path      string
	AppKey    string
	AppSecret string
}
