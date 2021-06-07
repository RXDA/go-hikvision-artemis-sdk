# golang hikvision artemis sdk

### import
```
go get -u github.com/RXDA/go-hikvision-artemis-sdk
```

### usage
```golang
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
```