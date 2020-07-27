package WatchHandler

import "atlas-service/pkg/api/zk"

type ConfigHandler struct {
}

func (c *ConfigHandler) nodeCreate(path zk.EventNodeCreated) {

}
func (c *ConfigHandler) nodeDeleted(path zk.EventNodeDeleted) {

}
func (c *ConfigHandler) nodeDataChanged(path zk.EventNodeDataChanged) {

}
func (c *ConfigHandler) nodeChildrenChanged(path zk.EventNodeChildrenChanged) {

}
