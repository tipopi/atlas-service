package zk

type BaseEvent struct {
	path   string
	client *Client
}

//用于eventBus识别观察者传递数据
type EventNodeCreated struct {
	BaseEvent
}
type EventNodeDeleted struct {
	BaseEvent
}
type EventNodeDataChanged struct {
	BaseEvent
}
type EventNodeChildrenChanged struct {
	BaseEvent
}
