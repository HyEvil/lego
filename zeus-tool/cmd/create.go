package cmd

import (
	"io/ioutil"

	"github.com/flosch/pongo2"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"go.yym.plus/zeus-tool/tpl"
	"go.yym.plus/zeus-tool/utils"
)

var (
	createCmd = &cobra.Command{
		Use:   "create [type]",
		Short: "create zeus bean",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			beanType := args[0]
			beanName := args[1]
			switch beanType {
			case "model":
				if err := generateModel(beanName); err != nil {
					cmd.PrintErr(err)
				}
			case "service":
				if err := generateService(beanName); err != nil {
					cmd.PrintErr(err)
				}
			case "repository":
				if err := generateRepository(beanName); err != nil {
					cmd.PrintErr(err)
				}
			case "bll":
				if err := generateBLL(beanName); err != nil {
					cmd.PrintErr(err)
				}
			default:
				cmd.PrintErr(errors.Errorf("not support create %s", beanType))
			}
		},
	}
	beanFile        *string
	repositoryTable *string
	repositoryModel *string
	forceWrite      *bool
)

func init() {
	rootCmd.AddCommand(createCmd)
	beanFile = createCmd.PersistentFlags().String("file", "", "")
	repositoryTable = createCmd.PersistentFlags().String("table", "", "")
	repositoryModel = createCmd.PersistentFlags().String("model", "", "")
	forceWrite = createCmd.PersistentFlags().Bool("force", false, "")
}

func generateModel(beanName string) error {
	tpl :=
		`package model

	type {{ bean }} struct {
		Id         int64 {{ idTag |safe}}
	}
`
	fileName := *beanFile
	if fileName == "" {
		fileName = strcase.ToLowerCamel(beanName)
	}
	fileName = fileName + ".go"
	path := utils.ProjectFilePathEx(true, "dao", "model", fileName)
	ctx := pongo2.Context{
		"bean":  strcase.ToCamel(beanName),
		"idTag": "`xorm:\"pk autoincr\"`",
	}
	return generateBeanAndWriteFile(path, tpl, ctx)
}

func generateService(beanName string) error {
	tpl :=
		`package service
func {{ beanConstruct}}() *{{beanStruct}} {
	return &{{beanStruct}}{}
}

type {{beanStruct}} struct {
}
`
	fileName := *beanFile
	if fileName == "" {
		fileName = strcase.ToLowerCamel(beanName)
	}
	fileName = fileName + ".go"
	path := utils.ProjectFilePathEx(true, "service", fileName)
	ctx := pongo2.Context{
		"beanConstruct": strcase.ToCamel(beanName) + "Server",
		"beanStruct":    strcase.ToLowerCamel(beanName) + "Server",
	}
	return generateBeanAndWriteFile(path, tpl, ctx)
}

func generateBLL(beanName string) error {
	tpl :=
		`package bll
var (
	{{beanVar}} = {{beanStruct}}{}
)

type {{beanStruct}} struct {
}`
	fileName := *beanFile
	if fileName == "" {
		fileName = strcase.ToLowerCamel(beanName)
	}
	fileName = fileName + ".go"
	path := utils.ProjectFilePathEx(true, "dao", "bll", fileName)
	ctx := pongo2.Context{
		"beanVar":    strcase.ToCamel(beanName),
		"beanStruct": strcase.ToLowerCamel(beanName),
	}
	return generateBeanAndWriteFile(path, tpl, ctx)
}

func generateRepository(beanName string) error {
	fileName := *beanFile
	if fileName == "" {
		fileName = strcase.ToLowerCamel(beanName)
	}
	tableName := *repositoryTable
	if tableName == "" {
		tableName = beanName
	}
	modelName := *repositoryModel
	if modelName == "" {
		modelName = beanName
	}
	fileName = fileName + ".go"
	path := utils.ProjectFilePathEx(true, "dao", "repository", fileName)
	ctx := pongo2.Context{
		"beanVar":    strcase.ToCamel(beanName),
		"beanStruct": strcase.ToLowerCamel(beanName),
		"table":      strcase.ToCamel(tableName),
		"model":      strcase.ToCamel(modelName),
	}
	return generateBeanAndWriteFile(path, tpl.Repository, ctx)
}

func generateBeanAndWriteFile(fileName string, tplContent string, ctx pongo2.Context) error {
	if utils.FileExists(fileName) && !*forceWrite {
		return errors.New("has exist")
	}
	tpl := pongo2.Must(pongo2.FromString(tplContent))
	data, err := tpl.ExecuteBytes(ctx)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		return err
	}
	return utils.GoFormat(fileName)
}
