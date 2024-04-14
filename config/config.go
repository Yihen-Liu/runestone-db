// Package config /**
package config

import (
	"bytes"
	"encoding/json"
	"os"
)

var Version = "dev-dirty"

type Log struct {
	LogLevel   int    `json:"LogLevel"`
	LogFileDir string `json:"LogFileDir"`
}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Port     string
	Db       string
}

type AppConfig struct {
	Logger   Log               `json:"Logger"`
	ChainRpc map[string]string `json:"ChainRpc"`
	Database Database          `json:"Database"`
}

var AppConf AppConfig

func init() {
	file, err := os.ReadFile(ConfigName)
	if err != nil {
		panic("config  file error:" + err.Error())
	}
	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	config := AppConfig{}
	if err := json.Unmarshal(file, &config); err != nil {
		panic("unmarshal json config err:" + err.Error())
	}
	AppConf = config
}
