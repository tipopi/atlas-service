package zk

import "strings"

var config *Client

func SetConfig(hosts string, path string) ([]byte, error) {
	var err error
	config, err = NewWithHandler(strings.Split(hosts, ","), &ConfigHandler{})
	if err != nil {
		return nil, err
	}
	err = config.EventRegistry()
	if err != nil {
		return nil, err
	}

	config.Watch(path)
	//config.Get(path)
	return config.Get(path)

}

func GetConfig() *Client {
	return config
}
