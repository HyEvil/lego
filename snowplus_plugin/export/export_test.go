package export

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"testing"
)

func TestGenerateStructTransform2(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()
	jsonData := `
{"Root":"/Users/yym/go/src/code.snowplus.cn/grpc-server/sfa","SourcePkg":"code.snowplus.cn/pb/sfa","SourceTypeName":"","SourceFileName":"/Users/yym/go/pkg/mod/code.snowplus.cn/pb/sfa@v0.0.0-20200215091755-d06a715062fb/salesman.pb.go","SourceTypeLine":108,"SourceTypeLineOff":12,"TargetPkg":"code.snowplus.cn/grpc-server/sfa/model","TargetTypeName":"","TargetFileName":"/Users/yym/go/src/code.snowplus.cn/grpc-server/sfa/model/salesman.go","TargetTypeLine":46,"TargetTypeLineOff":14,"VarName":"req","GoExecutableDir":"/usr/local/go/bin"}
`
	config := GenerateStructConfig{
	}
	json.Unmarshal([]byte(jsonData), &config)
	_, err := GenerateStructAssignCodeEx(&config)
	fmt.Print(123)
	fmt.Println(err)
	if err != nil {
		t.Error(err)
	}
}
