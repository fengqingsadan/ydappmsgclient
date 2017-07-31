package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"cindasoft.com/library/utils"
)

var (
	host   string
	ydHttp utils.JGHttp

	GetAccessTokenUrl = "http://%s/cgi/gettoken"
	SendMsgUrl        = "http://%s/cgi/msg/send?accessToken=%s"
)

//----------
type AppInfo struct {
	Buin      int32  `json:"buin"`
	Host      string `json:"host"`
	AppId     string `json:"appId"`
	AppAesKey string `json:"appAesKey"`

	KeyBytes []byte `json:"keyBytes"`
}

func (i *AppInfo) Valid() bool {
	yes := i.Buin > 0
	if yes {
		yes = i.Host != ""
	}
	if yes {
		yes = i.AppId != ""
	}
	if yes {
		yes = i.AppAesKey != ""
	}
	if yes {
		var err error
		i.KeyBytes, err = base64.StdEncoding.DecodeString(i.AppAesKey)
		yes = (err == nil)
	}
	return yes
}

func (i *AppInfo) String() string {
	return fmt.Sprintf("buin: %d \r host:%s \r appId:%s \r appAesKey:%s", i.Buin, i.Host, i.AppId, i.AppAesKey)
}

//----------
type SysMsg struct {
	Title string        `json:"title"`
	Msg   []interface{} `json:"msg"`
}

func (m *SysMsg) Valid() bool {
	yes := m.Title != ""
	if yes {
		yes = m.Msg != nil && len(m.Msg) > 0
	}
	return yes
}

func (m *SysMsg) Data() ([]byte, error) {
	return json.Marshal(m)
}

//----------
type SendTo struct {
	ToUser string `json:"toUser"`
	ToDept string `json:"toDept"`
}

func (t *SendTo) Valid() bool {
	return t.ToUser != "" || t.ToDept != ""
}

//----------
type JsonConfig struct {
	App *AppInfo `json:"appInfo"`
	Msg *SysMsg  `json:"sysMsg"`
	To  *SendTo  `json:"sendTo"`
}

func (c *JsonConfig) Valid() bool {
	yes1 := c.App != nil && c.App.Valid()
	yes2 := c.Msg != nil && c.Msg.Valid()
	yes3 := c.To != nil && c.To.Valid()
	return yes1 && yes2 && yes3
}

//----------
type AccessTokenInfo struct {
	Token    string `json:"accessToken"`
	ExpireIn int64  `json:"expireIn"`
}

type Text struct {
	Content string `json:"content"`
}

type Link struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	Action int32  `json:"action,omitempty"`
}

//----------
type AppMsgBase struct {
	ToUser  string `json:"toUser"`
	ToDept  string `json:"toDept"`
	MsgType string `json:"msgType"`
}

//----------

type AppText struct {
	Text *Text `json:"text"`
}

//----------
type AppLink struct {
	Link *Link `json:"link"`
}

//----------
type AppSysMsg struct {
	AppMsgBase
	Msg *SysMsg `json:"sysMsg"`
}

//----------
type AccessTokenResult struct {
	Code    int    `json:"errcode"`
	Msg     string `json:"errmsg"`
	Encrypt string `json:"encrypt"`
	Data    []byte `json:"data,omitempty"`
}
