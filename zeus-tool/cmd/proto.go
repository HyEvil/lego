package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"go.yym.plus/zeus-tool/proto"
	"go.yym.plus/zeus-tool/utils"
)

const (
	apiDirName = "api"
)

var (
	protoCmd = &cobra.Command{
		Use:   "proto [action]",
		Short: "zeus proto operation",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			setPath()
			action := ""
			if len(args) == 0 {
				action = "gen"
			} else {
				action = args[0]
			}
			switch action {
			case "gen":
				if err := generate(); err != nil {
					cmd.PrintErr(err)
				}

			case "install":
				if err := checkProtoc(true); err != nil {
					cmd.PrintErr(err)
				}
			}

		},
	}
	pbFormat   *bool
	genSwagger *bool
	genHttp    *bool
)

func init() {
	rootCmd.AddCommand(protoCmd)
	pbFormat = protoCmd.PersistentFlags().Bool("format", true, "")
	genSwagger = protoCmd.PersistentFlags().Bool("swagger", true, "")
	genHttp = protoCmd.PersistentFlags().Bool("http", true, "")
}

func setPath() {
	path := os.Getenv("PATH")
	//path = path + string(os.PathListSeparator) + filepath.Join(utils.GoPath(), "src")
	path = path + string(os.PathListSeparator) + filepath.Join(utils.GoPath(), "bin")
	os.Setenv("PATH", path)
}

func generate() error {
	if err := generateProto(); err != nil {
		return errors.WithMessage(err, "gen proto error")
	}

	if err := generateGoProto(); err != nil {
		return errors.WithMessage(err, "gen go error")
	}

	if *genSwagger {
		if err := generateSwagger(); err != nil {
			return errors.WithMessage(err, "gen swagger error")
		}
	}

	return nil
}

func generateProto() error {
	if err := checkProtoTools(false); err != nil {
		return err
	}

	protoDirPath := utils.ProjectPath(apiDirName)

	if !utils.FileExists(protoDirPath) {
		return errors.New("proto dir not exist")
	}
	outDirPath := utils.ProjectDirPathEx(true, apiDirName, "generated")
	if !utils.FileExists(outDirPath) {
		os.MkdirAll(outDirPath, 0666)
	}

	protoFiles, err := filepath.Glob(filepath.Join(protoDirPath, "*.zeus"))
	if err != nil {
		return err
	}

	pbPaths := []string{}
	converter := proto.NewZeusConverter()
	for _, protoPath := range protoFiles {
		baseName := filepath.Base(protoPath)
		fileName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		pbPath := filepath.Join(outDirPath, fileName+".proto")

		err = converter.Parse(protoPath)
		if err != nil {
			return errors.WithMessagef(err, "parser file :%s", protoPath)
		}
		pbContent, _ := converter.PBContent()
		err = ioutil.WriteFile(pbPath, []byte(pbContent), 0666)
		if err != nil {
			return errors.WithMessagef(err, "write pb file :%s", pbPath)
		}
		if *genHttp {
			httpContent, _ := converter.HttpContent()
			if err != nil {
				return errors.WithMessagef(err, "gen http error :%s", protoPath)
			}
			httpFilePath := filepath.Join(outDirPath, fileName+".http.go")
			err = ioutil.WriteFile(httpFilePath, []byte(httpContent), 0666)
			if err != nil {
				return errors.WithMessagef(err, "write http file :%s", pbPath)
			}
			utils.GoFormat(httpFilePath)
		}

		pbPaths = append(pbPaths, pbPath)
	}
	if *pbFormat {
		err = utils.RunCmd("clang-format -i " + strings.Join(pbPaths, " "))
		if err != nil {
			return errors.WithMessagef(err, "format error")
		}
	}
	return nil
}

func checkProtoTools(install bool) error {
	if err := checkProtoc(install); err != nil {
		return err
	}
	if err := checkGoGoProto(install); err != nil {
		return err
	}

	if err := checkGenSwagger(install); err != nil {
		return err
	}
	if err := checkClangFormat(install); err != nil {
		return err
	}
	if err := checkPackr(install); err != nil {
		return err
	}
	return nil
}

func checkClangFormat(install bool) error {
	if _, err := exec.LookPath("clang-format"); err != nil {
		if !install {
			return errors.Errorf("clang-format not exist,please run install")
		}
		switch runtime.GOOS {
		case "darwin":
			return utils.RunCmd("brew install clang-format")
		case "linux":
			return utils.RunCmd("snap install --classic clang-format")
		default:
			return errors.New("您还没安装clang-format，请进行手动安装")
		}
	}
	return nil
}

func checkPackr(install bool) error {
	if _, err := exec.LookPath("packr2"); err != nil {
		if !install {
			return errors.Errorf("packr2 not exist,please run install")
		}
		return utils.RunCmd("go get -u github.com/gobuffalo/packr/v2/packr2")
	}
	return nil
}

func checkProtoc(install bool) error {
	if _, err := exec.LookPath("protoc"); err != nil {
		if !install {
			return errors.Errorf("google protobuf not exist,please run install")
		}
		switch runtime.GOOS {
		case "darwin":
			return utils.RunCmd("brew install protobuf")
		case "linux":
			return utils.RunCmd("snap install --classic protobuf")
		default:
			return errors.New("您还没安装protobuf，请进行手动安装：https://github.com/protocolbuffers/protobuf/releases")
		}
	}
	return nil
}

func checkGoGoProto(install bool) error {
	if _, err := exec.LookPath("protoc-gen-gofast"); err != nil {
		if !install {
			return errors.Errorf("protoc-gen-gofast not exist")
		}
		if err := utils.RunCmd("go get -u github.com/gogo/protobuf/protoc-gen-gofast"); err != nil {
			return err
		}
	}
	return nil
}

func checkGenSwagger(install bool) error {
	if _, err := exec.LookPath("protoc-gen-bswagger"); err != nil {
		if !install {
			return errors.Errorf("protoc-gen-bswagger not exist")
		}
		if err := utils.RunCmd("go get -u github.com/go-kratos/kratos/tool/protobuf/protoc-gen-bswagger"); err != nil {
			return err
		}
	}
	return nil
}

func generateGoProto() error {
	pbProtoPath := utils.ProjectPath(apiDirName, "generated")
	//outPath := utils.ProjectPath(apiDirName)
	protoFiles, err := filepath.Glob(filepath.Join(pbProtoPath, "*.proto"))
	if err != nil {
		return err
	}
	if len(protoFiles) == 0 {
		return errors.New("no proto files")
	}

	cmd := fmt.Sprintf("protoc --proto_path=%s --proto_path=%s --proto_path=%s --gofast_out=plugins=grpc:%s",
		pbProtoPath,
		filepath.Join(utils.GoPath(), "src"),
		filepath.Join(utils.GoPath(), "src", "go.yym.plus", "zeus-tool", "protobuf"),
		pbProtoPath)
	for _, file := range protoFiles {
		cmd += " " + file
	}
	return utils.RunCmd(cmd)
}

func generateSwagger() error {
	protoPath := utils.ProjectPath(apiDirName, "generated")
	outPath := utils.ProjectDirPathEx(true, apiDirName, "generated", "swagger")
	utils.DeleteFileByPattern(filepath.Join(outPath, "*.json"))
	protoFiles, err := filepath.Glob(filepath.Join(protoPath, "*.proto"))
	if err != nil {
		return err
	}
	if len(protoFiles) == 0 {
		return errors.New("no proto files")
	}

	cmd := fmt.Sprintf("protoc --proto_path=%s --proto_path=%s --proto_path=%s --bswagger_out=:%s",
		protoPath,
		filepath.Join(utils.GoPath(), "src"),
		filepath.Join(utils.GoPath(), "src", "go.yym.plus", "zeus-tool", "protobuf"),
		outPath)

	outFiles := []string{}
	for _, file := range protoFiles {
		cmd += " " + file
		outFileName := strings.TrimSuffix(filepath.Base(file), ".proto")
		outFiles = append(outFiles, filepath.Join(outPath, fmt.Sprintf("%s.swagger.json", outFileName)))
	}

	err = utils.RunCmd(cmd)
	if err != nil {
		return err
	}

	mergedContent, err := mergeSwagger(outFiles)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(outPath, "swagger.json"), mergedContent, 0666)
	if err != nil {
		return err
	}

	if *genHttp {
		err := utils.RunCmd("packr2", protoPath)
		if err != nil {
			return err
		}
		var registerTpl = pongo2.Must(pongo2.FromString(
	`package proto
		import (
		nethttp "net/http"
		ginSwagger "github.com/swaggo/gin-swagger"
		)
		func RegisterSwagger(swaggerPath string, httpServer *http.Engine) {
			box := packr.New("swagger", "./swagger")
			httpServer.StaticFS("/static/swagger/", box)
			url := ginSwagger.URL("/static/swagger/swagger.json")
			httpServer.GET(path.Join(swaggerPath, "/*any"), ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		}`))
		content, err := registerTpl.Execute(pongo2.Context{

		})
		if err != nil {
			return err
		}
		buff := bytes.Buffer{}
		//buff.WriteString(fmt.Sprintf("package proto\n"))
		buff.WriteString(content)
		err = ioutil.WriteFile(path.Join(protoPath, "swagger.http.go"), buff.Bytes(), 0666)
		if err != nil {
			return err
		}
		utils.GoFormat(path.Join(protoPath, "swagger.http.go"))

	}
	return nil
}

func mergeSwagger(files []string) ([]byte, error) {
	sw := struct {
		Swagger string `json:"swagger"`
		Info    struct {
			Title   string `json:"title"`
			Version string `json:"version"`
		} `json:"info"`
		Schemes     []string               `json:"schemes"`
		Consumes    []string               `json:"consumes"`
		Produces    []string               `json:"produces"`
		Paths       map[string]interface{} `json:"paths"`
		Definitions map[string]interface{} `json:"definitions"`
	}{}

	sw.Swagger = "2.0"
	sw.Info.Title = "api.proto"
	sw.Info.Version = ""
	sw.Schemes = []string{"http", "https"}
	sw.Consumes = []string{
		"application/json",
		"multipart/form-data"}
	sw.Produces = []string{"application/json"}
	sw.Paths = map[string]interface{}{}
	sw.Definitions = map[string]interface{}{}
	for _, file := range files {
		input, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		inputData := struct {
			Paths       map[string]interface{}
			Definitions map[string]interface{}
		}{}
		err = json.Unmarshal(input, &inputData)
		if err != nil {
			return nil, err
		}
		for k, v := range inputData.Paths {
			sw.Paths[k] = v
		}

		for k, v := range inputData.Definitions {
			sw.Definitions[k] = v
		}
	}
	return json.MarshalIndent(&sw, "", "	")
}
