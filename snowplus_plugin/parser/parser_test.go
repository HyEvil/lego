package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	typ1, _ := ParseTypeFromPackage("../test/gendao/model", "ChannelTypea")
	//	typ2, _ := ParseTypeFromPackage("../test/gendao/model", "ChannelTypeaa")
	fmt.Println(typ1.String())
}
