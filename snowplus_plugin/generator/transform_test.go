package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"yym/snowpluslib/parser"
	"yym/snowpluslib/utils"
)

func TestGenerateStructTransform(t *testing.T) {
	sourceTypeName := "Test1"
	targetTypeName := "Test2"
	typ1, err := parser.ParseTypeFromPackage("yym/snowpluslib/test", "F:\\work\\go\\src\\yym\\snowpluslib", sourceTypeName)
	if err != nil {
		panic(err)
	}
	//typ2, err := parser.ParseTypeFromPackage("code.snowplus.cn/pb/sfa", "SalesmanStoreDetailResp")
	typ2, err := parser.ParseTypeFromPackage("yym/snowpluslib/test", "F:\\work\\go\\src\\yym\\snowpluslib", targetTypeName)
	if err != nil {
		panic(err)
	}

	generated, err := GenerateStructTransform(true, typ1.Underlying().AsStructWrapper(), typ2.Underlying().AsStructWrapper(), "aaa", nil, parser.MakeFieldNotContainsFilter("XXX_"))
	if err != nil {
		panic(err)
	}
	file := "../test/testgen.go"
	if utils.FileExists(file) {
		err := os.Remove(file)
		if err != nil {
			panic(err)
		}
	}
	data := fmt.Sprintf(`package test
func init(){
aaa:=%s{}
b:=%s{%s}
fmt.Println(b)
}`, sourceTypeName, targetTypeName, generated)
	err = ioutil.WriteFile("../test/testgen.go", []byte(data), 0666)
	if err != nil {
		panic(err)
	}
	err = utils.GoFormat("../test/testgen.go")
	if err != nil {
		panic(err)
	}
}
