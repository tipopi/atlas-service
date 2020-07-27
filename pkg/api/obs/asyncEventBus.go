package obs

type AsyncEventBus struct {
	EventBus //组合
}

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
