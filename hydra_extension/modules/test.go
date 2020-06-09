package modules

import (
	"fmt"
	"yym/hydra_extension/hydra"
)

func init() {

	hydra.RegisterType("TestType", func() (*Test, error) {
		return &Test{}, nil
	})
	hydra.RegisterModule("Test", &Test{})
}

type Test struct {
	a int
}

type info struct {
	A int
	B string
}

func (self *Test) Test1(a* string, b int,c []string,  d info) (string, int) {
	return "hello", 1
}

func (self* Test) Test2(a map[string]string) {
	return
}

func (self Test) Test3(a []interface{}) (Test) {
	return Test{a:1}
}

func (self* Test) Test4( d hydra.Duration) {
	fmt.Println(d)
	return
}