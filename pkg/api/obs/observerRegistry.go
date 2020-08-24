package obs

import (
	"atlas-service/pkg/api/log"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type ObserverRegistry struct {
	//map并发读写会panic，加个读写锁压压惊，之后重构为cas或cow
	sync.RWMutex
	m map[reflect.Type][]ObserverAction
}

func (o *ObserverRegistry) registry(observer interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				err = r.(error)
			case string:
				err = errors.New(r.(string))
			default:
				err = errors.New("unknown exception")
			}
			log.Error(err.Error())
		}
	}()
	if o.m == nil {
		o.m = make(map[reflect.Type][]ObserverAction)
	}
	for eventType, eventAction := range buildActions(observer) {
		var actions []ObserverAction
		actions = o.m[eventType]
		var exit bool
		o.Lock()
		if _, exit = o.m[eventType]; !exit {
			a := make([]ObserverAction, 0, 10)
			o.m[eventType] = a
		}
		o.m[eventType] = append(actions, eventAction...)
		o.Unlock()
	}
	return err
}

//向外部提供此event下的actions
func (o *ObserverRegistry) getObserverActions(event interface{}) []ObserverAction {
	o.RLock()
	defer o.RUnlock()
	//精准匹配，暂时不考虑类型组合
	return o.m[reflect.TypeOf(event)]
}

//以event Type为键组装action
func buildActions(observer interface{}) map[reflect.Type][]ObserverAction {
	actionMap := make(map[reflect.Type][]ObserverAction)
	t := reflect.TypeOf(observer)
	for _, method := range getMethod(observer) {
		//获取方法第一个参数
		eventType := method.Type().In(0)
		var actions []ObserverAction
		actions = actionMap[eventType]
		var exit bool
		if _, exit = actionMap[eventType]; !exit {
			a := make([]ObserverAction, 0, 10)
			actionMap[eventType] = a
		}
		actionMap[eventType] = append(actions, ObserverAction{t, method})
	}
	return actionMap
}

//获取observer定义的所有订阅方法
func getMethod(observer interface{}) []reflect.Value {
	v := reflect.ValueOf(observer)
	//判断是函数

	fv := v
	switch v.Kind() {
	case reflect.Func:
		return []reflect.Value{v}
	case reflect.Ptr:
		fv = v.Elem()
	case reflect.Struct:
	default:
		panic("Type error: Please input struct or function")
	}
	//默认订阅
	m := v.MethodByName("Subscribe")
	methods := make([]reflect.Value, 0, 10)
	if m.IsValid() {
		methods = append(methods, m)
	}

	//获取tag上的自定义订阅 todo:找不到字段会panic
	if filed, exit := fv.Type().FieldByName("SubscribeMethod"); exit {
		if tags, tagExit := filed.Tag.Lookup("methods"); tagExit {
			tagValues := strings.Split(tags, ",")
			for _, value := range tagValues {
				method := v.MethodByName(value)
				if !method.IsValid() {
					log.Error(fmt.Sprintf("The observer does not have this method:%s", value))
					continue
				}
				methods = append(methods, method)
			}
		}
	}
	return methods

}
