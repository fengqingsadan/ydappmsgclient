package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"cindasoft.com/library/slog"
	"cindasoft.com/library/utils"
)

type AppClient struct {
	sync.RWMutex

	Buin              int32
	Host              string
	AppId             string
	AppAesKey         []byte
	AccessToken       string
	AccessTokenExpire int64

	ToUsers string
	ToDepts string
}

func NewAppClient(buin int32, host, appId, appAesKey string) (*AppClient, error) {
	key, err := base64.StdEncoding.DecodeString(appAesKey)
	if err != nil {
		return nil, err
	}
	c := &AppClient{
		Buin:      buin,
		AppId:     appId,
		AppAesKey: key,
	}
	return c, nil
}

//----------
func (c *AppClient) SetAppInfo(app *AppInfo) {
	c.Buin = app.Buin
	c.Host = app.Host
	c.AppId = app.AppId
	c.AppAesKey = app.KeyBytes
}

//----------
func (c *AppClient) SetSendTo(to *SendTo) {
	c.ToUsers = to.ToUser
	c.ToDepts = to.ToDept
}

//----------
func (c *AppClient) SendSysMsg(msg *SysMsg) {
	sysMsg := &AppSysMsg{}
	sysMsg.ToUser = c.ToUsers
	sysMsg.ToDept = c.ToDepts
	sysMsg.MsgType = "sysMsg"
	sysMsg.Msg = msg

	data, err := json.Marshal(sysMsg)
	if err != nil {
		slog.Exit("[appclient] send sys-msg error: ", err)
	}

	enMsg, err := AesEncrypt(data, c.AppAesKey, c.AppId)
	if err != nil {
		slog.Error("[appclient] send sys-msg error: aes encrypt failed: ", err)
		return
	}

	m := make(map[string]interface{})
	m["buin"] = c.Buin
	m["appId"] = c.AppId
	m["encrypt"] = enMsg
	data, _ = json.Marshal(m)
	token := c.getAccessToken()
	msgSendUrl := fmt.Sprintf(SendMsgUrl, c.Host, token)
	data, st := ydHttp.DoTextPost(msgSendUrl, string(data))
	if !st.IsStatusOK() {
		slog.Error("[appclient] send sys-msg error: post request error:", err)
		return
	}

	slog.Info("[appclient] send sys-msg done, get result:", string(data))
}

//----------
func (c *AppClient) getAccessToken() string {
	now := utils.CurrentTimeSecond()
	c.Lock()
	if c.AccessTokenExpire-now > 5*60 {
		c.Unlock()
		return c.AccessToken
	}
	c._getAccessTokenFromYD()
	c.Unlock()
	return c.AccessToken
}

func (c *AppClient) _getAccessTokenFromYD() {
	timex := fmt.Sprintf("%d", utils.CurrentTimeMillis())
	enMsg, err := AesEncrypt([]byte(timex), c.AppAesKey, c.AppId)
	if err != nil {
		slog.Error("[appclient] get access-token error: aes encrypt failed: ", err)
		return
	}

	rm := make(map[string]interface{})
	rm["buin"] = c.Buin
	rm["appId"] = c.AppId
	rm["encrypt"] = enMsg
	bs, _ := json.Marshal(rm)
	httpUrl := fmt.Sprintf(GetAccessTokenUrl, c.Host)
	data, st := ydHttp.DoTextPost(httpUrl, string(bs))
	if !st.IsStatusOK() {
		slog.Error("[appclient] get access-token error: post request failed: ", st)
		return
	}

	var result AccessTokenResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		slog.Error("[appclient] get access-token error: unmarshal result failed: ", err)
		return
	}

	if result.Code != 0 {
		slog.Error("[appclient] get access-token error: get access_token result is not ok: ", string(data))
		return
	}

	rawMsg, err := AesDecrypt(result.Encrypt, c.AppAesKey)
	if err != nil {
		slog.Error("[appclient] get access-token error: aes decrypt error: ", err)
		return
	}

	var token AccessTokenInfo
	err = json.Unmarshal(rawMsg.Data, &token)
	if err != nil {
		slog.Error("[appclient] get access-token error: unmarshal decrypt result error: ", err, string(rawMsg.Data))
		return
	}

	c.AccessToken = token.Token
	c.AccessTokenExpire = token.ExpireIn + utils.CurrentTimeSecond()
	slog.Info("[appclient] get  access-token success: ", token)
}
