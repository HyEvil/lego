package hydra

import (
	"go.uber.org/atomic"
	"math"
	"sync"
)

var (
	instances      sync.Map
	nextInstanceId = atomic.NewUint32(0)
	instanceIdPool = sync.Pool{
		New: func() interface{} {
			nextInstanceId.CAS(math.MaxUint32, 0)
			return nextInstanceId.Inc()
		},
	}
)

func PutInstance(instance interface{}) uint32 {

	id := instanceIdPool.Get().(uint32)
	instances.Store(id, instance)
	return id
}

func RemoveInstance(id uint32) {
	instances.Delete(id)
	instanceIdPool.Put(id)
}

func GetInstance(id uint32) (interface{}, bool) {
	return instances.Load(id)
}

func init() {
	RegisterFunc("Instance.delete", func(id uint32) {
		RemoveInstance(id)
	})
}
