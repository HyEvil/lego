package export

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
	_ "yym/snowpluslib/log"
)

//export GenerateStructAssignCode
func GenerateStructAssignCode(configDataPtr unsafe.Pointer, configSize C.int) *C.char {

	configData := C.GoBytes(configDataPtr, configSize)
	config := GenerateStructConfig{}

	err := json.Unmarshal(configData, &config)
	if err != nil {
		//return C.CString(err.Error())
		logrus.WithError(err).Error("GenerateStructAssignCode error")
		return nil
	}
	ret, err := GenerateStructAssignCodeEx(&config)
	if err != nil {
		//return C.CString(err.Error())
		logrus.WithError(err).Error("GenerateStructAssignCode error")
		return nil
	}
	cstr := C.CString(ret)
	time.AfterFunc(time.Second*60, func() {
		C.free(unsafe.Pointer(cstr))
	})
	return cstr
}
