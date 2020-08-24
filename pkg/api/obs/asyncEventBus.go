package obs

import "reflect"

type AsyncEventBus struct {
	EventBus //组合
}

var asyncEventBus *AsyncEventBus

func (eventBus *AsyncEventBus) Post(event interface{}) {
	runFunc = asyncRun
	eventBus.EventBus.Post(event) //委托给同步eventbus,但实际run方法已被替换为异步
}

func (eventBus *AsyncEventBus) PostWithObs(event interface{}, obs interface{}) {
	runFunc = asyncRun
	eventBus.EventBus.PostWithObs(event, obs)
}

func asyncRun(action ObserverAction, event interface{}) {
	go action.execute(event)
}
func GetAsyncEventBus() *AsyncEventBus {
	if asyncEventBus == nil {
		asyncEventBus = &AsyncEventBus{EventBus{ObserverRegistry{m: make(map[reflect.Type][]ObserverAction)}}}
	}
	return asyncEventBus
}
