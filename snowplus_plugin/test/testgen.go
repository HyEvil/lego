package test

import (
	"fmt"
	"github.com/spf13/cast"
)

func init() {
	aaa := Test1{}
	b := Test2{
		V1: int32(aaa.V1),
		V2: NamedString(aaa.V2),
		V3: int(aaa.V3),
		V4: NamedString(aaa.V4),
		V5: &aaa.V5,
		V6: *aaa.V6,
		V7: func() *int {
			v := int(*aaa.V7)
			return &v
		}(),
		V8: func() *NamedString {
			v := NamedString(cast.ToString(int(**aaa.V8)))
			return &v
		}(),
		V11: func() []int8 {
			s := []int8{}
			for _, e := range aaa.V11 {
				s = append(s, int8(e))
			}
			return s
		}(),
		V12: func() []*NamedString {
			s := []*NamedString{}
			for _, e := range aaa.V12 {
				s = append(s, func() *NamedString {
					v := NamedString(*e)
					return &v
				}())
			}
			return s
		}(),
	}
	fmt.Println(b)
}
