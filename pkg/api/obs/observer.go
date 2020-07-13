package obs

type observer struct {
	//自定义观察者方法
	SubscribeMethod string  `methods:""`
}

//默认作为观察者的方法
func(o *observer) subscribe(event interface{})  {

}
