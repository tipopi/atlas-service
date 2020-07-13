package obs

type EventBus struct {
	ObserverRegistry
}

func (eventBus *EventBus)Register(observer interface{})  {
	eventBus.registry(observer)
}

func (eventBus *EventBus) Post(event interface{})  {
	for _,action:=range eventBus.getObserverActions(event){
		run(action,event)
	}

}
func (eventBus *EventBus) AsyncPost(event interface{})  {
	for _,action:=range eventBus.getObserverActions(event){
		go run(action,event)
	}
}
func run(action ObserverAction,event interface{})  {
	action.execute(event)
}