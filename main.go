package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"time"

	"cindasoft.com/library/slog"
)

func init() {
	flag.Set("logmode", "stdout:debug,file:debug")
	slog.InitByFlags()
	slog.Info("--------sysmsg client start-----------")
}

func main() {
	cfg := readConfig()
	client := &AppClient{}
	client.SetAppInfo(cfg.App)
	client.SetSendTo(cfg.To)
	client.SendSysMsg(cfg.Msg)
	select {
	case <-time.After(3 * time.Second):
		break
	}
}

func readConfig() *JsonConfig {
	data, err := ioutil.ReadFile("app.json")
	if err != nil {
		slog.Exit("read app.json failed: ", err)
	}

	var cfg JsonConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		slog.Exit("read app.json failed: ", string(data), err)
	}
	if !cfg.Valid() {
		slog.Exit("read app.json invalid")
	}

	return &cfg
}
