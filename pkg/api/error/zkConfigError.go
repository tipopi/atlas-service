package error

type ZkConfigError struct {
	Msg string
}

func (e ZkConfigError) Error() string {
	return e.Msg
}
