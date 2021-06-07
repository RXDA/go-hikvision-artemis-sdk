package hk_artemis_sdk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)



func initialBasicHeader(method, path string, headers, querys, bodys map[string]string, signHeaderPrefixList []string, appKey, appSecret string) map[string]string {
	headers["x-ca-timestamp"] = strconv.FormatInt(time.Now().UnixNano(), 10)
	headers["x-ca-nonce"] = uuid.New().String()
	headers["x-ca-key"] = appKey
	headers["x-ca-signature"] = sign(appSecret, method, path, headers, querys, bodys, signHeaderPrefixList)
	return headers
}

func sign(secret, method, path string, headers, querys, bodys map[string]string, signHeaderPrefixList []string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(buildStringToSign(method, path, headers, querys, bodys, signHeaderPrefixList)))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func buildStringToSign(method, path string, headers, querys, bodys map[string]string, signHeaderPrefixList []string) string {
	var sb strings.Builder
	sb.WriteString(strings.ToUpper(method))
	sb.WriteString("\n")
	if len(headers) != 0 {
		if val, ok := headers["Accept"]; ok {
			sb.WriteString(val)
			sb.WriteString("\n")
		}
		if val, ok := headers["Content-MD5"]; ok {
			sb.WriteString(val)
			sb.WriteString("\n")
		}
		if val, ok := headers["Content-Type"]; ok {
			sb.WriteString(val)
			sb.WriteString("\n")
		}
		if val, ok := headers["Date"]; ok {
			sb.WriteString(val)
			sb.WriteString("\n")
		}
	}
	sb.WriteString(buildHeaders(headers, signHeaderPrefixList))
	sb.WriteString(buildResource(path, querys, bodys))
	return sb.String()
}

func buildResource(path string, querys, bodys map[string]string) string {
	sb := strings.Builder{}
	if !isBlankString(path) {
		sb.WriteString(path)
	}

	sortMap := treemap.NewWithStringComparator()
	if len(querys) != 0 {
		for k, v := range querys {
			if !isBlankString(k) {
				sortMap.Put(k, v)
			}
		}
	}

	if len(bodys) != 0 {
		for k, v := range querys {
			if !isBlankString(k) {
				sortMap.Put(k, v)
			}
		}
	}

	sbParam := strings.Builder{}
	it := sortMap.Iterator()
	for it.Next() {
		k, v := it.Key().(string), it.Value().(string)
		if !isBlankString(k) {
			if sbParam.Len() > 0 {
				sbParam.WriteString("&")
			}
			sbParam.WriteString(k)
			if !isBlankString(v) {
				sbParam.WriteString("=")
				sbParam.WriteString(v)
			}
		}
	}

	if sbParam.Len() > 0 {
		sb.WriteString("?")
		sb.WriteString(sbParam.String())
	}
	return sb.String()
}

func buildHeaders(headers map[string]string, signHeaderPrefixList []string) string {
	var sb strings.Builder
	if len(signHeaderPrefixList) != 0 {
		//remove x-ca-signature, Accept, Content-MD5,Content-Type,Date
		for i := len(signHeaderPrefixList) - 1; i >= 0; i-- {
			if signHeaderPrefixList[i] == "x-ca-signature" || signHeaderPrefixList[i] == "Accept" ||
				signHeaderPrefixList[i] == "Content-MD5" || signHeaderPrefixList[i] == "Content-Type" ||
				signHeaderPrefixList[i] == "Date" {
				signHeaderPrefixList = append(signHeaderPrefixList[:i], signHeaderPrefixList[i+1:]...)
			}
		}
		sort.Strings(signHeaderPrefixList)
	}

	if len(headers) != 0 {
		m := mapToTreeMap(headers)
		signHeadersStringBuilder := strings.Builder{}

		it := m.Iterator()
		for it.Next() {
			k, v := it.Key(), it.Value()
			sk := k.(string)
			sv := v.(string)
			if isHeaderToSign(sk, signHeaderPrefixList) {
				sb.WriteString(sk)
				sb.WriteString(":")

				if !isBlankString(sv) {
					sb.WriteString(sv)
				}
				sb.WriteString("\n")
				if sb.Len() > 0 {
					signHeadersStringBuilder.WriteString(",")
				}
				signHeadersStringBuilder.WriteString(sk)
			}
		}
		headers["x-ca-signature-headers"] = signHeadersStringBuilder.String()
	}
	return sb.String()
}

func initURL(host, path string, querys map[string]string) string {
	vals := url.Values{}
	for k, v := range querys {
		vals.Add(k, v)
	}

	u := url.URL{
		Path:     path,
		RawQuery: vals.Encode(),
	}
	if strings.HasPrefix(host, "https") {
		u.Scheme = "https"
		u.Host = strings.TrimPrefix(host, "https://")
	} else {
		u.Scheme = "http"
		u.Host = strings.TrimPrefix(host, "http://")
	}

	return u.String()
}

var client = http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// HttpGet http get method
func (aReq ArtemisReq) HttpGet(headers, querys map[string]string, signHeaderPrefixList []string, timeout time.Duration) (*ArtemisResp, error) {
	headers = map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/text;charset=UTF-8",
	}

	headers = initialBasicHeader("GET", aReq.Path, headers, querys, nil, signHeaderPrefixList, aReq.AppKey, aReq.AppSecret)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", initURL(aReq.Host, aReq.Path, querys), nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := ArtemisResp{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
