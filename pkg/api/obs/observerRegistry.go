package obs

import (
	"reflect"
	"strings"
	"sync"
)

type ObserverRegistry struct {
	//map并发读写会崩溃，加个读写锁压压惊，之后重构为cas或cow
	sync.RWMutex
	m map[reflect.Type][]ObserverAction
}

func(o *ObserverRegistry) registry(observer interface{})  {
	for eventType,eventAction:=range buildActions(observer){
		var actions []ObserverAction
		var exit bool
		if actions,exit=o.m[eventType];!exit{
			a:=make([]ObserverAction,0,10)
			o.Lock()
			o.m[eventType]=a
			o.Unlock()
			actions=a
		}
		actions=append(actions,eventAction...)
	}
}
//向外部提供此event下的actions
func (o *ObserverRegistry)getObserverActions(event interface{}) []ObserverAction {
	o.RLock()
	defer o.Unlock()
	//精准匹配，暂时不考虑类型组合
	return o.m[reflect.TypeOf(event)]
}

//以event Type为键组装action
func buildActions(observer interface{})map[reflect.Type][]ObserverAction{
	actionMap:=make(map[reflect.Type][]ObserverAction)
	t:=reflect.TypeOf(observer)
	for _,method:=range getMethod(observer){
		//获取方法第一个参数
		eventType:=method.Type().In(0)
		var actions []ObserverAction
		var exit bool
		if actions,exit=actionMap[eventType];!exit{
			a:=make([]ObserverAction,0,10)
			actionMap[eventType]=a
			actions=a
		}
		actions=append(actions,ObserverAction{t,method})

	}
	return actionMap
}
//获取observer定义的所有订阅方法
func getMethod(observer interface{}) []reflect.Value  {
	v:=reflect.ValueOf(observer)
	m:=v.MethodByName("subscribe")
	methods:=make([]reflect.Value,0,10)
	if !m.IsNil(){
		methods=append(methods,m)
	}

	if filed,exit:=v.Type().FieldByName("subscribeMethod");exit{
		if tags,tagExit:=filed.Tag.Lookup("methods");tagExit{
			tagValues:=strings.Split(tags,",")
			for _,value:=range tagValues{
				method:=v.MethodByName(value)
				if method.IsNil() {
					//has error
					continue
				}
				methods = append(methods,method )
			}
		}
	}
	return methods

}