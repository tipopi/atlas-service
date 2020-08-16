package error

import "atlas-service/pkg/api/log"

func CheckError(e error, useLog bool, f func(err error)) {
	if e != nil {
		if useLog {
			log.Error(e.Error())
		}
		f(e)
	}
}
