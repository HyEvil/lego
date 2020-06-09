package export

import "C"
import (
	"encoding/json"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"yym/snowpluslib/generator"
	"yym/snowpluslib/parser"
)

var (
	hasSetEnv = false
)

type GenerateStructConfig struct {
	Root              string
	SourcePkg         string
	SourceTypeName    string
	SourceFileName    string
	SourceTypeLine    int
	SourceTypeLineOff int
	TargetPkg         string
	TargetTypeName    string
	TargetFileName    string
	TargetTypeLine    int
	TargetTypeLineOff int
	VarName           string
	GoExecutableDir   string
}

func GenerateStructAssignCodeEx(cfg *GenerateStructConfig) (string, error) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.WithField("stack", string(debug.Stack())).Error("GenerateStructAssignCodeEx panic")
		}
	}()
	jsonData, _ := json.Marshal(cfg)
	logrus.WithField("json", string(jsonData)).Debug("GenerateStructAssignCode json data")
	if !hasSetEnv {
		hasSetEnv = true
		if runtime.GOOS == "darwin" {
			path := os.Getenv("PATH")
			path = path + ":" + cfg.GoExecutableDir
			os.Setenv("PATH", path)
		}
	}
	sourceType, err := parser.ParseTypeByFileLineOffset(cfg.SourcePkg, cfg.Root, cfg.SourceFileName, cfg.SourceTypeLine, cfg.SourceTypeLineOff)
	if err != nil {
		return "", err
	}
	targetType, err := parser.ParseTypeByFileLineOffset(cfg.TargetPkg, cfg.Root, cfg.TargetFileName, cfg.TargetTypeLine, cfg.TargetTypeLineOff)
	if err != nil {
		return "", err
	}
	return generator.GenerateStructTransform(cfg.SourcePkg == cfg.TargetPkg, sourceType.AsStructWrapper(), targetType.AsStructWrapper(),
		cfg.VarName, parser.MakeFieldNotContainsFilter("XXX_"),
		parser.MakeFieldNotContainsFilter("XXX_"))
}
