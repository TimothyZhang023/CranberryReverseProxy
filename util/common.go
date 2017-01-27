package util

import (
	"fmt"
	"reflect"
	"runtime"
	"errors"
	"os"

	"github.com/zts1993/crp/log"

)

func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Zeus() {
	if e := recover(); e != nil {

		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]

		log.Errorf("[panic] err=%v \n %s \n ", e, buf)
		fmt.Printf("[panic] err=%v \n %s \n ", e, buf)

		os.Exit(-1)
	}
}

func PanicTest(info string) {
	panic(info)
}

func PanicProcess(function interface{}, callback interface{}) (result []reflect.Value, err error) {

	if err := recover(); err != nil {

		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]

		log.Errorf("[%s] panic:%v %s \n", GetFuncName(function), err, buf)
		//fmt.Printf("[%s] panic:%v %s \n", GetFuncName(function), err, buf)

		if callback != nil {
			return ReflectCall(callback)
		}
	}

	return
}

func ReflectCall(m interface{}, params ...interface{}) (result []reflect.Value, err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[%s] panic:%v", GetFuncName(ReflectCall), err)
		}
	}()

	f := reflect.ValueOf(m)
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
