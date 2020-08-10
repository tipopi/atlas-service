package zk

import (
	"atlas-service/pkg/api/zk/Handler"
	"strings"
)

var (
	config *Client
)

func SetConfig(hosts string, path string) ([]byte, error) {

	config, err := NewWithHandler(strings.Split(hosts, ","), &Handler.ConfigHandler{})
	if err != nil {
		return nil, err
	}
	return config.Get(path)

}

func GetConfig() *Client {
	return config
}
