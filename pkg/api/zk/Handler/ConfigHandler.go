package Handler

import "atlas-service/pkg/api/zk"

type ConfigHandler struct {
}

func (c *ConfigHandler) NodeCreate(path zk.EventNodeCreated) {

}
func (c *ConfigHandler) NodeDeleted(path zk.EventNodeDeleted) {

}
func (c *ConfigHandler) NodeDataChanged(path zk.EventNodeDataChanged) {

}
func (c *ConfigHandler) NodeChildrenChanged(path zk.EventNodeChildrenChanged) {

}
