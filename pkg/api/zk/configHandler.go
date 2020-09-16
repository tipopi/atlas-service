package zk

import (
	error2 "atlas-service/pkg/api/error"
	"atlas-service/pkg/api/log"
	"bytes"
	"github.com/spf13/viper"
)

type ConfigHandler struct {
	SubscribeMethod string `methods:"NodeDataChanged"`
}

func (c *ConfigHandler) NodeCreate(event EventNodeCreated) {

}
func (c *ConfigHandler) NodeDeleted(event EventNodeDeleted) {

}
func (c *ConfigHandler) NodeDataChanged(event EventNodeDataChanged) {
	viper.SetConfigType("yaml")
	config, err := event.client.Get(event.path)
	error2.CheckError(err, true, func(err error) {
		log.Error("config get error")
	})
	err = viper.ReadConfig(bytes.NewBuffer(config))
	error2.CheckError(err, true, func(err error) {
		log.Error("config reload error")
	})
	log.Info("test:" + viper.GetString("test"))
}
func (c *ConfigHandler) NodeChildrenChanged(event EventNodeChildrenChanged) {

}
func (c *ConfigHandler) Subscribe(event EventNodeDataChanged) {
	log.Info("test" + string(event.path))
}
