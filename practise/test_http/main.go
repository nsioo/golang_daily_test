package main

import (
	"bytes"
	"code.byted.org/gopkg/logs"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type DeviceRegister struct {
	UserUniqueId string `json:"user_unique_id"`
	AppId        uint32 `json:"app_id"`
	Os           string `json:"os"`
}

type deviceRegisterResponse struct {
	DeviceId     uint64 `json:"device_id"`
	InstallId    uint64 `json:"install_id"`
	BdDid        string `json:"bd_did"`
	Cd           string `json:"cd"`
	InstallIdStr string `json:"install_id_str"`
	NewUser      uint8  `json:"new_user"`
	Ssid         string `json:"ssid"`
	ServerTime   uint64 `json:"server_time"`
}

// 根据 user_unique_id 和 app_id 注册 device_id
// 仅适用于私有化
func (dr DeviceRegister) RegisterDeviceId() (string, error) {
	bodyJson, err := json.Marshal(dr.generateBody())
	if err != nil {
		logs.Error("marshal body err: %v", err)
		return "", err
	}

	logs.Warn("req body is %v", dr.generateBody())

	req, err := http.NewRequest("POST", "http://10.225.130.116/service/2/device_register/", bytes.NewBuffer(bodyJson))
	if err != nil {
		logs.Error("new request err: %v", err)
		return "", err
	}
	//req.Header.Set("User-Agent", "Data Creator/2.0.0 (OnPremise)")
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("app_id", strconv.Itoa(int(dr.AppId)))
	//req.Header.Set("Host", "snssdk.vpc.com")
	req.Host = "snssdk.vpc.com"
	host := req.Header.Get("Host")
	fmt.Println(host)

	client := &http.Client{Timeout: time.Second * 30}
	times := 2
	var resp *http.Response
	for {
		if times <= 0 {
			break
		}
		resp, err = client.Do(req)
		if err != nil {
			logs.Error("http upload err: %+v", err.Error())
		} else {
			break
		}

		times--
	}

	if resp == nil {
		logs.Error("resp is nil")
		return "", errors.New("resp is nil")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read resp body err: %v", err)
		return "", err
	}

	res := deviceRegisterResponse{}
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		logs.Error("unmarshal resp err: %v", err)
		logs.Warn("res is %v", string(respBody))
		time.Sleep(time.Second)
		return "", err
	}

	if res.DeviceId == 0 && res.BdDid == "" && res.Cd == "" {
		return "", errors.New("generate device_id err")
	}

	logs.Warn("bd_did is %v", res.BdDid)

	return res.BdDid, nil
}

func (dr DeviceRegister) generateBody() map[string]interface{} {
	header := make(map[string]interface{})

	body := make(map[string]interface{}, 0)

	body["aid"] = dr.AppId
	body["user_unique_id"] = dr.UserUniqueId

	// os 需要区分 ios 和 android，对应枚举值 iOS \ ANDROID
	body["os"] = formatOs(dr.Os)

	// ios 需要填 vendor_id， android 需要填 openudid
	uniqueIdr, _ := uuid.GenerateUUID()
	if dr.Os == "ios" {
		body["vendor_id"] = strings.ToUpper(uniqueIdr)
	} else {
		body["openudid"] = strings.ToUpper(uniqueIdr)
	}

	header["header"] = body

	return header
}

func formatOs(osName string) string {
	if osName == "ios" {
		return "iOS"
	}

	return "ANDROID"
}

func main() {
	dr := DeviceRegister{
		UserUniqueId: "276095447832965",
		AppId:        10000012,
		Os:           "ios",
	}

	did, err := dr.RegisterDeviceId()
	if err != nil {
		logs.Error("err: %v", err)
		return
	}

	fmt.Println(did)
}
