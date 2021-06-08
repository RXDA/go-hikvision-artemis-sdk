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
	Schema    string // https or http
	Host      string // only host
	Port      uint16 // port
	Path      string
	AppKey    string
	AppSecret string
}
