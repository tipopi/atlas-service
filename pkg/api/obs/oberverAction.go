package obs

import "reflect"

type ObserverAction struct {
	target interface{}
	method reflect.Value
}

func (o *ObserverAction) execute(event interface{}) {
	input := make([]reflect.Value, 1)
	input = append(input, reflect.ValueOf(event))
	o.method.Call(input)
}
