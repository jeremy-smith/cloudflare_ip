package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	ApiListZones       = "https://api.cloudflare.com/client/v4/zones"
	ApiListDNSRecords  = "https://api.cloudflare.com/client/v4/zones/%s/dns_records"
	ApiCreateDNSRecord = "https://api.cloudflare.com/client/v4/zones/%s/dns_records"
	ApiUpdateDNSRecord = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s"
)

type Cloudflare struct {
	bearerToken string
	client      http.Client
}

type successError struct {
	Success bool                     `json:"success"`
	Errors  []map[string]interface{} `json:"errors"`
}

type jsonResultRecord struct {
	Id string `json:"id"`
}

type jsonResultSingle struct {
	successError
	Result jsonResultRecord
}

type jsonResultMultiple struct {
	successError
	Result []jsonResultRecord
}

func PrettyPrintJSON(b []byte) {
	bB := bytes.Buffer{}
	_ = json.Indent(&bB, b, "", "  ")
	fmt.Println(bB.String())
}

// compileErrStr takes an error array from Cloudflare and returns the errors messages as a string
func compileErrStr(errs []map[string]interface{}) string {
	errStr := ""
	for i, e := range errs {
		v, ok := e["message"].(string)
		if ok {
			errStr += v
			if i < len(errs)-1 {
				errStr += ", "
			}
		}
	}
	return errStr
}

func NewCloudflare(bearerToken string) Cloudflare {
	c := http.Client{}
	return Cloudflare{
		bearerToken: bearerToken,
		client:      c,
	}
}

func (cf Cloudflare) ListZones(domain string) (string, error) {
	return _cfGet(&cf, domain, "", "", "")
}

func (cf Cloudflare) ListDNSRecords(zoneID, recordName, recordType string) (string, error) {
	return _cfGet(&cf, "", zoneID, recordName, recordType)
}

func _cfGet(cf *Cloudflare, domain, zoneID, recordName, recordType string) (string, error) {
	var apiPath string

	queryParams := make(map[string]string)

	if domain != "" { // List Zones
		apiPath = ApiListZones
		queryParams["name"] = domain
	} else { // List DNS Records
		apiPath = fmt.Sprintf(ApiListDNSRecords, zoneID)
		queryParams["name"] = recordName
		queryParams["type"] = recordType
	}

	req, err := http.NewRequest(http.MethodGet, apiPath, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", "Bearer "+cf.bearerToken)
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := cf.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jr jsonResultMultiple
	_ = json.Unmarshal(b, &jr)

	if !jr.Success {
		errStr := compileErrStr(jr.Errors)
		return "", errors.New("Cloudflare returned an error: " + errStr)

	}

	if len(jr.Result) == 0 {
		return "", nil
	}

	return jr.Result[0].Id, nil
}

func (cf Cloudflare) CreateDNSRecord(zoneID, domain, recordType, ip string) (string, error) {
	return _cfUpdateCreate(&cf, zoneID, "", domain, recordType, ip)
}

func (cf Cloudflare) UpdateDNSRecord(zoneID, dnsID, domain, recordType, ip string) (string, error) {
	return _cfUpdateCreate(&cf, zoneID, dnsID, domain, recordType, ip)
}

func _cfUpdateCreate(cf *Cloudflare, zoneID, dnsID, domain, recordType, ip string) (string, error) {
	var method string
	var apiPath string

	if dnsID == "" {
		method = http.MethodPost
		apiPath = fmt.Sprintf(ApiCreateDNSRecord, zoneID)
	} else {
		method = http.MethodPut
		apiPath = fmt.Sprintf(ApiUpdateDNSRecord, zoneID, dnsID)
	}

	postData := []byte(fmt.Sprintf(
		`{"type":"%s", "name":"%s", "content":"%s"}`, recordType, domain, ip))

	body := bytes.NewBuffer(postData)

	req, err := http.NewRequest(method, apiPath, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", "Bearer "+cf.bearerToken)

	resp, err := cf.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jr jsonResultSingle
	_ = json.Unmarshal(b, &jr)

	if !jr.Success {
		errStr := compileErrStr(jr.Errors)
		return "", errors.New("Cloudflare returned an error: " + errStr)
	}

	return jr.Result.Id, nil
}
