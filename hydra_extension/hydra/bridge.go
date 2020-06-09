package hydra

/*
#include "hydra.h"
*/
import "C"

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
	"yym/hydra_extension/3rd/codec"
)

var (
	hydraResume C.HydraResume
	hydraRet    C.HydraRet
)

type MsgpackCoder struct {
	encoder *codec.Encoder
	decoder *codec.Decoder
}

var (
	codecHandler codec.MsgpackHandle
	coderPool    = sync.Pool{
		New: func() interface{} {
			return &MsgpackCoder{
				encoder: codec.NewEncoderBytes(nil, &codecHandler),
				decoder: codec.NewDecoderBytes(nil, &codecHandler),
			}
		},}
)

//export HydraInit
func HydraInit(f1 C.HydraResume, f2 C.HydraRet) {
	hydraResume = f1
	hydraRet = f2
}

func init() {
	codecHandler.MapType = reflect.TypeOf(map[string]interface{}{})
}

//export SuspendCall
func SuspendCall(goEx unsafe.Pointer, coroutineId C.uint, namePtr *C.char, msgData unsafe.Pointer, msgSize C.int) {
	cmd := C.GoString(namePtr)
	handler := extensionHandlers[cmd]
	if handler == nil {
		resumeError(goEx, coroutineId, fmt.Errorf("cmd %s not have handler", cmd))
		return
	}
	argData := C.GoBytes(msgData, msgSize)
	go onSuspendCall(goEx, coroutineId, argData, handler)
}

//export Call
func Call(goEx unsafe.Pointer, namePtr *C.char, msgData unsafe.Pointer, msgSize C.int) C.char {
	cmd := C.GoString(namePtr)
	handler := extensionHandlers[cmd]
	if handler == nil {
		hydraReturnError(goEx, fmt.Errorf("cmd %s not have handler", cmd))
		return C.char(0)
	}
	argData := C.GoBytes(msgData, msgSize)
	var retData []byte
	coder := coderPool.Get().(*MsgpackCoder)
	decoder := coder.decoder
	decoder.ResetBytes(argData)

	encoder := coder.encoder
	encoder.ResetBytes(&retData)
	err := handler(encoder, decoder)
	if err != nil {
		hydraReturnError(goEx, err)
		coderPool.Put(coder)
		return C.char(0)
	} else {
		hydraReturn(goEx, retData)
		coderPool.Put(coder)
		return C.char(1)
	}
}

func onSuspendCall(goEx unsafe.Pointer, coroutineId C.uint, argData []byte, handler func(encoder *codec.Encoder, decoder *codec.Decoder) error) {
	coder := coderPool.Get().(*MsgpackCoder)
	decoder := coder.decoder
	decoder.ResetBytes(argData)
	var retData []byte
	encoder := coder.encoder
	encoder.ResetBytes(&retData)
	err := handler(encoder, decoder)
	if err != nil {
		resumeError(goEx, coroutineId, err)
	} else {
		resume(goEx, coroutineId, retData)
	}
	coderPool.Put(coder)
}

func resumeError(goEx unsafe.Pointer, coroutineId C.uint, err error) {
	what := []byte(err.Error())
	C.CallHydraResume(hydraResume, goEx, coroutineId, C.char(0), unsafe.Pointer(&what[0]), C.uint(len(what)))
}

func resume(goEx unsafe.Pointer, coroutineId C.uint, data []byte) {
	if len(data) == 0 {
		C.CallHydraResume(hydraResume, goEx, coroutineId, C.char(1), unsafe.Pointer(nil), C.uint(0))
	} else {
		C.CallHydraResume(hydraResume, goEx, coroutineId, C.char(1), unsafe.Pointer(&data[0]), C.uint(len(data)))
	}
}

func hydraReturn(goEx unsafe.Pointer, data []byte) {
	if len(data) == 0 {
		C.CallHydraRet(hydraRet, goEx, unsafe.Pointer(nil), C.uint(0))
	} else {
		C.CallHydraRet(hydraRet, goEx, unsafe.Pointer(&data[0]), C.uint(len(data)))
	}

}

func hydraReturnError(goEx unsafe.Pointer, err error) {
	what := []byte(err.Error())
	C.CallHydraRet(hydraRet, goEx, unsafe.Pointer(&what[0]), C.uint(len(what)))
}
