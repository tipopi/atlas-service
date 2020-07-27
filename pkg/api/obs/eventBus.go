package obs

import "reflect"

type EventBus struct {
	ObserverRegistry
}

var runFunc = run

func (eventBus *EventBus) Register(observer interface{}) {
	eventBus.registry(observer)
}

func (eventBus *EventBus) Post(event interface{}) {
	for _, action := range eventBus.getObserverActions(event) {
		runFunc(action, event)
	}
}

//只有observer类型和event类型都符合才会执行
func (eventBus *EventBus) PostWithObs(event interface{}, obs interface{}) {
	t := reflect.TypeOf(obs)
	for _, action := range eventBus.getObserverActions(event) {
		if reflect.TypeOf(action.target) == t {
			runFunc(action, event)
		}
	}
}

func run(action ObserverAction, event interface{}) {
	action.execute(event)
}
