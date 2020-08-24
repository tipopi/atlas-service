package obs

import "reflect"

type ObserverAction struct {
	target interface{}
	method reflect.Value
}

func (o *ObserverAction) execute(event interface{}) {
	input := []reflect.Value{reflect.ValueOf(event)}
	o.method.Call(input)
}
