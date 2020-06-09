package parser

import "sync"

var (
	cache sync.Map
)

func GetCachedType(key, packageName, name string) *TypeWrapper {
	v, ok := cache.Load(key)
	if ok {
		return v.(*TypeWrapper)
	}
	var err error
	v, err = ParseTypeFromPackage(packageName, "", name)
	if err != nil {
		return nil
	}
	return v.(*TypeWrapper)
}
