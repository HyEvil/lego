package hydra

/*
#include "hydra.h"
*/
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"reflect"
	"strings"
	"yym/hydra_extension/3rd/codec"
)

var (
	extensionHandlers = make(map[string]func(encoder *codec.Encoder, decoder *codec.Decoder) error)
	registeredTypes   = map[reflect.Type]interface{}{}
)

func RegisterModule(name string, m interface{}) {
	module := reflect.TypeOf(m)
	if module.Kind() != reflect.Ptr {
		panic("must be as struct  pointer")
	}
	instance := reflect.ValueOf(m)
	for i := 0; i < module.NumMethod(); i++ {
		method := module.Method(i)

		fixName := method.Name[:1]
		if strings.ToUpper(fixName) != fixName {
			continue
		}
		if len(method.Name) > 1 {
			fixName = strings.ToLower(fixName) + method.Name[1:]
		}
		methodValue := instance.MethodByName(method.Name)
		RegisterFunc(name+"."+fixName, methodValue.Interface())
	}
}

func RegisterType(name string, creator interface{}) {
	creatorType := reflect.TypeOf(creator)
	var module reflect.Type
	if creatorType.Kind() == reflect.Func {
		if creatorType.NumOut() != 2 || creatorType.Out(0).Kind() != reflect.Ptr {
			panic("crate func invalid")
		}
		method, ok := creatorType.Out(1).MethodByName("Error")
		if !ok {
			panic("crate func invalid")
		}
		if method.Type.NumOut() != 1 || method.Type.NumIn() != 0 || method.Type.Out(0).Kind() != reflect.String {
			panic("crate func invalid")
		}
		RegisterFuncEx(name+".new", creator, nil, func(values *[]reflect.Value) error {
			retValues := *values
			if retValues[1].IsNil() {
				id := PutInstance(retValues[0].Interface())
				*values = []reflect.Value{reflect.ValueOf(id)}
			} else {
				return errors.WithMessage(retValues[1].Interface().(error), "create instance failed")
			}
			return nil
		})
		module = creatorType.Out(0)
	} else {
		module = reflect.TypeOf(creator)
	}
	if module.Kind() != reflect.Ptr {
		panic("must be as struct  pointer")
	}
	registeredTypes[module.Elem()] = struct{}{}
	for i := 0; i < module.NumMethod(); i++ {
		method := module.Method(i)
		methodValue := method.Func
		//methodType:=method.Type
		fixName := method.Name[:1]
		if strings.ToUpper(fixName) != fixName {
			continue
		}
		if len(method.Name) > 1 {
			fixName = strings.ToLower(fixName) + method.Name[1:]
		}

		RegisterFuncEx(name+"."+fixName, methodValue.Interface(), func(argList []reflect.Value, argTypes []reflect.Type, decoder *codec.Decoder) (int, error) {
			var instanceId int
			err := decoder.Decode(&instanceId)
			if err != nil {
				return 0, fmt.Errorf("decode instanceId failed")
			}

			instance, ok := GetInstance(uint32(instanceId))
			if !ok || instance == nil {
				return 0, fmt.Errorf("instance not found")
			}
			argList[0] = reflect.ValueOf(instance)
			return 1, nil
		}, nil)
	}
}

func RegisterFuncEx(name string, f interface{}, before func(argList []reflect.Value, argTypes []reflect.Type, decoder *codec.Decoder) (int, error), after func(retValues *[]reflect.Value) error) {
	funcType := reflect.TypeOf(f)
	funcValue := reflect.ValueOf(f)
	if funcType.Kind() != reflect.Func {
		panic("must register with a create func")
	}
	argTypeCount := funcType.NumIn()
	argTypes := make([]reflect.Type, funcType.NumIn())
	for n := 0; n < argTypeCount; n++ {
		argTypes[n] = funcType.In(n)
	}
	retCount := funcType.NumOut()
	retNeedConvertToInstanceId := make([]*bool, retCount)
	retTypes := make([]reflect.Type, retCount)
	for n := 0; n < retCount; n++ {
		retTypes[n] = funcType.Out(n)
	}
	for i := 0; i < retCount; i++ {
		outType := funcType.Out(i)
		isPointer := false
		if outType.Kind() == reflect.Ptr {
			outType = outType.Elem()
			isPointer = true
		}
		if registeredTypes[outType] != nil {
			retNeedConvertToInstanceId[i] = &isPointer
		}
	}

	hasVariadic := funcType.IsVariadic()
	registerHandler(name, func(encoder *codec.Encoder, decoder *codec.Decoder) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("call panic:%v", r)
			}
		}()
		count := 0
		minArgCount := argTypeCount
		if hasVariadic {
			minArgCount = minArgCount - 1
		}
		argList := make([]reflect.Value, minArgCount)
		if before != nil {
			count, err = before(argList, argTypes, decoder)
			if err != nil {
				return err
			}
		}
		for i := count; i < argTypeCount; i++ {

			if hasVariadic && argTypeCount == i+1 {
				for {
					arg := reflect.New(argTypes[i].Elem())
					intf := arg.Interface()
					err = decoder.Decode(intf)
					if err != nil {
						if err == io.EOF {
							break
						}
						return errors.Errorf("arg %d not match:%v", i+1, err)
					}
					argList = append(argList, arg.Elem())
				}
			} else {
				arg := reflect.New(argTypes[i])
				intf := arg.Interface()
				err = decoder.Decode(intf)
				if err != nil {
					return errors.Errorf("arg %d not match:%v", i+1, err)
				}
				argList[i] = arg.Elem()
			}
			/*if argTypeCount == i+1 && hasVariadic {
				if err != nil {
					if err != io.EOF {
						return errors.Errorf("arg %d not match:%v", i+1, err)
					}
				}else{
					elem := arg.Elem()
					for n := 0; n < elem.Len(); n++ {
						argList = append(argList, elem.Index(n))
					}
				}

			} else {
				if err != nil {
					return errors.Errorf("arg %d not match:%v", i+1, err)
				}
				argList[i] = arg.Elem()
			}*/
		}

		rets := funcValue.Call(argList)
		if after != nil {
			if err := after(&rets); err != nil {
				return err
			}
		}
		for index, value := range rets {
			if retNeedConvertToInstanceId[index] != nil {
				var instance interface{}
				if *retNeedConvertToInstanceId[index] {
					instance = value.Interface()
				} else {
					newValue := reflect.New(retTypes[index])
					newValue.Elem().Set(value)
					instance = newValue.Interface()
				}
				id := PutInstance(instance)
				EncodeRet(encoder, id)
			} else {
				if err := EncodeRet(encoder, value.Interface()); err != nil {
					return errors.WithMessagef(err, "can not encode ret %d", index+1)
				}
			}

		}
		return nil
	})
}

func EncodeRet(encoder *codec.Encoder, v interface{}) error {
	if v != nil {
		if mayErr, ok := v.(error); ok {
			return encoder.Encode(mayErr.Error())
		} else {
			return encoder.Encode(v)
		}
	} else {
		return encoder.Encode(nil)
	}
}

func RegisterFunc(name string, f interface{}) {
	RegisterFuncEx(name, f, nil, nil)
}

func registerHandler(name string, handlder func(encoder *codec.Encoder, decoder *codec.Decoder) error) {
	extensionHandlers[name] = handlder
}

func TestFunc(name string, args ...interface{}) {
	handler := extensionHandlers[name]
	if handler == nil {
		panic("not handler")
	}
	en := codec.NewEncoderBytes(nil, &codecHandler)
	var data []byte = make([]byte, 0)
	en.ResetBytes(&data)
	for i := 0; i < len(args); i++ {
		en.Encode(args[i])
	}
	var tmp []byte
	en.ResetBytes(&tmp)
	de := codec.NewDecoderBytes(nil, &codecHandler)
	de.ResetBytes(data)

	err := handler(en, de)
	if err != nil {
		panic(err)
	}
}
