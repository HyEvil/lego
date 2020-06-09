package proto

import (
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/emicklei/proto"
	"github.com/emicklei/proto-contrib/pkg/protofmt"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cast"

	"go.yym.plus/zeus-tool/proto/parser"
	"go.yym.plus/zeus-tool/utils"
)

type StringBuffer struct {
	strings.Builder
}

type ServiceContext struct {
	resource       string
	middlewareList []string
}

type RpcContext struct {
	httpMethod     string
	httpPath       string
	middlewareList []string
}

type MessageContext struct {
	index      int
	preMessage *MessageContext
}

type FieldContext struct {
	jsonTag      string
	validatorTag string
	defaultTag   string
	nullable     string
	customTag    string
}

type ZeusProtoParser struct {
	lastError   error
	pbBuffer    StringBuffer
	httpContext *HttpContext
	parsedFile  map[string]interface{}
	ctx         struct {
		curService *ServiceContext
		curRpc     *RpcContext
		curMessage *MessageContext
		curField   *FieldContext
	}
}

type ZeusConverter struct {
	parser *ZeusProtoParser
}

func (self *ZeusProtoParser) EnterMapType(c *parser.MapTypeContext) {

}

func (self *ZeusProtoParser) ExitMapType(c *parser.MapTypeContext) {

}

func (self *ZeusProtoParser) EnterAnnotationArg(c *parser.AnnotationArgContext) {

}

func (self *ZeusProtoParser) ExitAnnotationArg(c *parser.AnnotationArgContext) {

}

func (self *ZeusProtoParser) EnterFieldType(c *parser.FieldTypeContext) {

}

func (self *ZeusProtoParser) ExitFieldType(c *parser.FieldTypeContext) {

}

func (self *ZeusProtoParser) EnterMessageDefineName(c *parser.MessageDefineNameContext) {
	//if _, ok := c.GetParent().(*parser.MessageContext); ok {
	self.pbBuffer.WriteSpacedAndLinkBreak(c.GetText(), "{")
	//}
}

func (self *ZeusProtoParser) ExitMessageDefineName(c *parser.MessageDefineNameContext) {

}

func (self *ZeusProtoParser) EnterRpcRequest(c *parser.RpcRequestContext) {
	req := c.GetText()
	if req == "" {
		req = "google.protobuf.Empty"
	}
	self.pbBuffer.WriteText("(", req, ") ")
}

func (self *ZeusProtoParser) EnterRpcResponse(c *parser.RpcResponseContext) {
	resp := c.GetText()
	if resp == "" {
		resp = "google.protobuf.Empty"
	}
	self.pbBuffer.WriteText("returns (", resp, ")")
	httpPath := self.ctx.curRpc.httpPath
	if self.ctx.curService.resource != "" {
		httpPath = path.Join(self.ctx.curService.resource, httpPath)
	}
	if self.ctx.curRpc.httpMethod != "" {
		self.pbBuffer.WriteTextAndLineBreak(" {")
		self.pbBuffer.WriteString(fmt.Sprintf(`option (google.api.http) = {
          %s: %s
        };
}`,
			self.ctx.curRpc.httpMethod,
			strconv.Quote(httpPath)))
	} else {
		self.pbBuffer.WriteEnding()
	}

}

func (self *ZeusProtoParser) ExitRpcRequest(c *parser.RpcRequestContext) {

}

func (self *ZeusProtoParser) ExitRpcResponse(c *parser.RpcResponseContext) {

}

func (self *ZeusProtoParser) VisitTerminal(node antlr.TerminalNode) {

}

func (self *ZeusProtoParser) VisitErrorNode(node antlr.ErrorNode) {
	fmt.Println(node)
}

func (self *ZeusProtoParser) EnterEveryRule(ctx antlr.ParserRuleContext) {

}

func (self *ZeusProtoParser) ExitEveryRule(ctx antlr.ParserRuleContext) {

}

func (self *ZeusProtoParser) EnterProto(c *parser.ProtoContext) {

}

func (self *ZeusProtoParser) EnterOptionList(c *parser.OptionListContext) {

}

func (self *ZeusProtoParser) EnterImportStatements(c *parser.ImportStatementsContext) {
	self.pbBuffer.WriteString(defaultPbImport())
	self.pbBuffer.WriteLineBreak()
}

func (self *ZeusProtoParser) EnterImportStatement(c *parser.ImportStatementContext) {
	importPath, _ := strconv.Unquote(c.StrLit().GetText())
	if strings.HasSuffix(importPath, ".zeus") {
		filePath := utils.LookPath(importPath)
		if filePath == "" {
			//	self.brea
		}
	} else {
		self.pbBuffer.WriteSpacedAndEnding("import", c.StrLit().GetText())
	}

	//	self.pbBuffer.WriteString()
}

func (self *ZeusProtoParser) EnterPackageStatement(c *parser.PackageStatementContext) {
	self.pbBuffer.WriteSpacedAndEnding(`syntax = "proto3"`)
	self.pbBuffer.WriteSpacedAndEnding("package", c.GetChild(1).(antlr.ParseTree).GetText())
	self.pbBuffer.WriteLineBreak()
	self.pbBuffer.WriteString(defaultOption())
	self.pbBuffer.WriteLineBreak()
}

func (self *ZeusProtoParser) EnterOption(c *parser.OptionContext) {

}

func (self *ZeusProtoParser) EnterOptionName(c *parser.OptionNameContext) {

}

func (self *ZeusProtoParser) EnterOptionBody(c *parser.OptionBodyContext) {

}

func (self *ZeusProtoParser) EnterOptionBodyVariable(c *parser.OptionBodyVariableContext) {

}

func (self *ZeusProtoParser) EnterTopLevelDef(c *parser.TopLevelDefContext) {

}

func (self *ZeusProtoParser) EnterMessage(c *parser.MessageContext) {
	curMessage := self.ctx.curMessage
	self.ctx.curMessage = &MessageContext{
		index:      1,
		preMessage: curMessage,
	}
	self.pbBuffer.WriteString("message ")
}

func (self *ZeusProtoParser) EnterMessageBody(c *parser.MessageBodyContext) {

}

func (self *ZeusProtoParser) EnterEnumDefinition(c *parser.EnumDefinitionContext) {
	self.pbBuffer.WriteString("enum " + c.EnumName().GetText() + c.EnumBody().GetText())
}

func (self *ZeusProtoParser) EnterEnumBody(c *parser.EnumBodyContext) {

}

func (self *ZeusProtoParser) EnterEnumField(c *parser.EnumFieldContext) {

}

func (self *ZeusProtoParser) EnterEnumValueOption(c *parser.EnumValueOptionContext) {

}

func (self *ZeusProtoParser) EnterExtend(c *parser.ExtendContext) {

}

func (self *ZeusProtoParser) EnterService(c *parser.ServiceContext) {
	self.ctx.curService = &ServiceContext{
	}
	self.pbBuffer.WriteString("service ")
}

func (self *ZeusProtoParser) EnterRpc(c *parser.RpcContext) {
	self.ctx.curRpc = &RpcContext{}
	self.pbBuffer.WriteString("rpc ")
}

func (self *ZeusProtoParser) EnterReserved(c *parser.ReservedContext) {

}

func (self *ZeusProtoParser) EnterRanges(c *parser.RangesContext) {

}

func (self *ZeusProtoParser) EnterRangeRule(c *parser.RangeRuleContext) {

}

func (self *ZeusProtoParser) EnterFieldNames(c *parser.FieldNamesContext) {

}

func (self *ZeusProtoParser) EnterTypeRule(c *parser.TypeRuleContext) {

}

func (self *ZeusProtoParser) EnterFieldNumber(c *parser.FieldNumberContext) {

}

func (self *ZeusProtoParser) EnterField(c *parser.FieldContext) {
	self.ctx.curField = &FieldContext{
		jsonTag:      "",
		validatorTag: "",
		defaultTag:   "",
		nullable:     "",
	}
}

func (self *ZeusProtoParser) EnterFieldOptions(c *parser.FieldOptionsContext) {

}

func (self *ZeusProtoParser) EnterFieldOption(c *parser.FieldOptionContext) {

}

func (self *ZeusProtoParser) EnterAnnotation(c *parser.AnnotationContext) {
	annotationName := c.GetText()
	annotationName = annotationName[1:]
	index := strings.Index(annotationName, "(")
	if index > 0 {
		annotationName = annotationName[:index]
	}
	t := AnnotationType(0)
	switch c.GetParent().(type) {
	case *parser.ServiceContext:
		t = AnnotationTypeService
	case *parser.RpcContext:
		t = AnnotationTypeRpc
	case *parser.MessageContext:
		t = AnnotationTypeMessage
	case *parser.FieldContext:
		t = AnnotationTypeField

	}
	var processor AnnotationProcessor = nil
	processors := messageAnnotationProcessors[t]
	if processors != nil {
		processor = processors[annotationName]
	}
	if processor == nil {
		panic(utils.WithPosError(c, "annotation %s not support", annotationName))
	} else {
		processor(self, c)
	}
}

func (self *ZeusProtoParser) EnterEmptyAnnotation(c *parser.EmptyAnnotationContext) {

}

func (self *ZeusProtoParser) EnterAnnotationWithArg(c *parser.AnnotationWithArgContext) {

}

func (self *ZeusProtoParser) EnterAnnotationWithNamedArg(c *parser.AnnotationWithNamedArgContext) {

}

func (self *ZeusProtoParser) EnterAnnotationName(c *parser.AnnotationNameContext) {

}

func (self *ZeusProtoParser) EnterAnnotationNamedArg(c *parser.AnnotationNamedArgContext) {

}

func (self *ZeusProtoParser) EnterAnnotationNamedArgKey(c *parser.AnnotationNamedArgKeyContext) {

}

func (self *ZeusProtoParser) EnterAnnotationNamedArgValue(c *parser.AnnotationNamedArgValueContext) {

}

func (self *ZeusProtoParser) EnterOneof(c *parser.OneofContext) {

}

func (self *ZeusProtoParser) EnterOneofField(c *parser.OneofFieldContext) {

}

func (self *ZeusProtoParser) EnterMapField(c *parser.MapFieldContext) {

}

func (self *ZeusProtoParser) EnterKeyType(c *parser.KeyTypeContext) {

}

func (self *ZeusProtoParser) EnterReservedWord(c *parser.ReservedWordContext) {

}

func (self *ZeusProtoParser) EnterFullIdent(c *parser.FullIdentContext) {

}

func (self *ZeusProtoParser) EnterMessageName(c *parser.MessageNameContext) {
}

func (self *ZeusProtoParser) EnterEnumName(c *parser.EnumNameContext) {

}

func (self *ZeusProtoParser) EnterMessageOrEnumName(c *parser.MessageOrEnumNameContext) {

}

func (self *ZeusProtoParser) EnterFieldName(c *parser.FieldNameContext) {

}

func (self *ZeusProtoParser) EnterOneofName(c *parser.OneofNameContext) {

}

func (self *ZeusProtoParser) EnterMapName(c *parser.MapNameContext) {

}

func (self *ZeusProtoParser) EnterServiceName(c *parser.ServiceNameContext) {
	self.pbBuffer.WriteSpacedAndLinkBreak(c.GetText(), "{")
	self.httpContext.NewService(c.GetText())
	self.httpContext.AddServiceMiddleware(self.ctx.curService.middlewareList...)
}

func (self *ZeusProtoParser) EnterRpcName(c *parser.RpcNameContext) {
	self.pbBuffer.WriteString(c.GetText())
	self.httpContext.NewHandler(c.GetText())
	self.httpContext.SetHandlerMethod(self.ctx.curRpc.httpMethod)
	self.httpContext.SetHandlerPath(path.Join(self.ctx.curService.resource, self.ctx.curRpc.httpPath))
	self.httpContext.AddHandlerMiddleware(self.ctx.curRpc.middlewareList...)
}

func (self *ZeusProtoParser) EnterMessageType(c *parser.MessageTypeContext) {

}

func (self *ZeusProtoParser) EnterMessageOrEnumType(c *parser.MessageOrEnumTypeContext) {

}

func (self *ZeusProtoParser) EnterEmptyStatement(c *parser.EmptyStatementContext) {

}

func (self *ZeusProtoParser) EnterConstant(c *parser.ConstantContext) {

}

func (self *ZeusProtoParser) ExitProto(c *parser.ProtoContext) {

}

func (self *ZeusProtoParser) ExitOptionList(c *parser.OptionListContext) {

}

func (self *ZeusProtoParser) ExitImportStatements(c *parser.ImportStatementsContext) {
	self.pbBuffer.WriteLineBreak()
}

func (self *ZeusProtoParser) ExitImportStatement(c *parser.ImportStatementContext) {

}

func (self *ZeusProtoParser) ExitPackageStatement(c *parser.PackageStatementContext) {

}

func (self *ZeusProtoParser) ExitOption(c *parser.OptionContext) {

}

func (self *ZeusProtoParser) ExitOptionName(c *parser.OptionNameContext) {

}

func (self *ZeusProtoParser) ExitOptionBody(c *parser.OptionBodyContext) {

}

func (self *ZeusProtoParser) ExitOptionBodyVariable(c *parser.OptionBodyVariableContext) {

}

func (self *ZeusProtoParser) ExitTopLevelDef(c *parser.TopLevelDefContext) {

}

func (self *ZeusProtoParser) ExitMessage(c *parser.MessageContext) {
	self.ctx.curMessage = self.ctx.curMessage.preMessage
	self.pbBuffer.WriteTextAndLineBreak("}")
}

func (self *ZeusProtoParser) ExitMessageBody(c *parser.MessageBodyContext) {

}

func (self *ZeusProtoParser) ExitEnumDefinition(c *parser.EnumDefinitionContext) {

}

func (self *ZeusProtoParser) ExitEnumBody(c *parser.EnumBodyContext) {

}

func (self *ZeusProtoParser) ExitEnumField(c *parser.EnumFieldContext) {

}

func (self *ZeusProtoParser) ExitEnumValueOption(c *parser.EnumValueOptionContext) {

}

func (self *ZeusProtoParser) ExitExtend(c *parser.ExtendContext) {

}

func (self *ZeusProtoParser) ExitService(c *parser.ServiceContext) {
	self.pbBuffer.WriteSpacedAndLinkBreak("}")
}

func (self *ZeusProtoParser) ExitRpc(c *parser.RpcContext) {
	self.ctx.curRpc = nil
}

func (self *ZeusProtoParser) ExitReserved(c *parser.ReservedContext) {
	self.ctx.curService = nil
}

func (self *ZeusProtoParser) ExitRanges(c *parser.RangesContext) {

}

func (self *ZeusProtoParser) ExitRangeRule(c *parser.RangeRuleContext) {

}

func (self *ZeusProtoParser) ExitFieldNames(c *parser.FieldNamesContext) {

}

func (self *ZeusProtoParser) ExitTypeRule(c *parser.TypeRuleContext) {

}

func (self *ZeusProtoParser) ExitFieldNumber(c *parser.FieldNumberContext) {

}

func (self *ZeusProtoParser) ExitField(c *parser.FieldContext) {
	fixFieldType := c.FieldType().GetText()
	if strings.Index(fixFieldType, "repeated") == 0 {
		fixFieldType = "repeated " + fixFieldType[len("repeated"):]
	}
	//	[(gogoproto.moretags) = 'validate:"omitempty,gte=1,lte=2" form:"gender"']
	self.pbBuffer.WriteSpaced(fixFieldType, c.FieldName().GetText(), "=", cast.ToString(self.ctx.curMessage.index), "[")
	moreTag := []string{}

	if self.ctx.curField.jsonTag == "" {
		self.ctx.curField.jsonTag = strconv.Quote(strcase.ToLowerCamel(c.FieldName().GetText()))
	}

	self.pbBuffer.WriteText("(gogoproto.jsontag) =", self.ctx.curField.jsonTag)

	if self.ctx.curField.customTag != "" {
		self.pbBuffer.WriteText(", (gogoproto.customtype) =", self.ctx.curField.customTag)
	}

	if self.ctx.curField.nullable != "" {
		self.pbBuffer.WriteText(", (gogoproto.nullable) =", self.ctx.curField.nullable)
	}

	if self.ctx.curField.defaultTag != "" {
		moreTag = append(moreTag, "default:"+self.ctx.curField.defaultTag)
	}

	if self.ctx.curField.validatorTag != "" {
		moreTag = append(moreTag, "validate:"+self.ctx.curField.validatorTag)
	}
	if len(moreTag) > 0 {
		self.pbBuffer.WriteText(", (gogoproto.moretags) = '")
		self.pbBuffer.WriteSpaced(moreTag[0], moreTag[1:]...)
		self.pbBuffer.WriteText("'")
	}

	self.pbBuffer.WriteTextAndEnding("]")
	self.ctx.curMessage.index++
	self.ctx.curField = nil
}

func (self *ZeusProtoParser) ExitFieldOptions(c *parser.FieldOptionsContext) {

}

func (self *ZeusProtoParser) ExitFieldOption(c *parser.FieldOptionContext) {

}

func (self *ZeusProtoParser) ExitAnnotation(c *parser.AnnotationContext) {

}

func (self *ZeusProtoParser) ExitEmptyAnnotation(c *parser.EmptyAnnotationContext) {

}

func (self *ZeusProtoParser) ExitAnnotationWithArg(c *parser.AnnotationWithArgContext) {

}

func (self *ZeusProtoParser) ExitAnnotationWithNamedArg(c *parser.AnnotationWithNamedArgContext) {

}

func (self *ZeusProtoParser) ExitAnnotationName(c *parser.AnnotationNameContext) {

}

func (self *ZeusProtoParser) ExitAnnotationNamedArg(c *parser.AnnotationNamedArgContext) {

}

func (self *ZeusProtoParser) ExitAnnotationNamedArgKey(c *parser.AnnotationNamedArgKeyContext) {

}

func (self *ZeusProtoParser) ExitAnnotationNamedArgValue(c *parser.AnnotationNamedArgValueContext) {

}

func (self *ZeusProtoParser) ExitOneof(c *parser.OneofContext) {

}

func (self *ZeusProtoParser) ExitOneofField(c *parser.OneofFieldContext) {

}

func (self *ZeusProtoParser) ExitMapField(c *parser.MapFieldContext) {

}

func (self *ZeusProtoParser) ExitKeyType(c *parser.KeyTypeContext) {

}

func (self *ZeusProtoParser) ExitReservedWord(c *parser.ReservedWordContext) {

}

func (self *ZeusProtoParser) ExitFullIdent(c *parser.FullIdentContext) {

}

func (self *ZeusProtoParser) ExitMessageName(c *parser.MessageNameContext) {
}

func (self *ZeusProtoParser) ExitEnumName(c *parser.EnumNameContext) {

}

func (self *ZeusProtoParser) ExitMessageOrEnumName(c *parser.MessageOrEnumNameContext) {

}

func (self *ZeusProtoParser) ExitFieldName(c *parser.FieldNameContext) {

}

func (self *ZeusProtoParser) ExitOneofName(c *parser.OneofNameContext) {

}

func (self *ZeusProtoParser) ExitMapName(c *parser.MapNameContext) {

}

func (self *ZeusProtoParser) ExitServiceName(c *parser.ServiceNameContext) {

}

func (self *ZeusProtoParser) ExitRpcName(c *parser.RpcNameContext) {

}

func (self *ZeusProtoParser) ExitMessageType(c *parser.MessageTypeContext) {

}

func (self *ZeusProtoParser) ExitMessageOrEnumType(c *parser.MessageOrEnumTypeContext) {

}

func (self *ZeusProtoParser) ExitEmptyStatement(c *parser.EmptyStatementContext) {

}

func (self *ZeusProtoParser) ExitConstant(c *parser.ConstantContext) {

}

func (self *StringBuffer) WriteSpaced(text string, after ...string) {
	self.WriteString(text)
	for _, s := range after {
		self.WriteString(" ")
		self.WriteString(s)
	}
}

func (self *StringBuffer) WriteText(text ...string) {
	for _, s := range text {
		self.WriteString(s)
	}
}

func (self *StringBuffer) WriteTextAndLineBreak(text ...string) {
	for _, s := range text {
		self.WriteString(s)
	}
	self.WriteLineBreak()
}

func (self *StringBuffer) WriteTextAndEnding(text ...string) {
	for _, s := range text {
		self.WriteString(s)
	}
	self.WriteEnding()
}

func (self *StringBuffer) WriteLineBreak() {
	self.WriteString("\r\n")
}

func (self *StringBuffer) WriteSpacedAndLinkBreak(text string, after ...string) {
	self.WriteString(text)
	for _, s := range after {
		self.WriteString(" ")
		self.WriteString(s)
	}
	self.WriteLineBreak()
}

func (self *StringBuffer) WriteSpacedAndEnding(text string, after ...string) {
	self.WriteString(text)
	for _, s := range after {
		self.WriteString(" ")
		self.WriteString(s)
	}
	self.WriteEnding()
}

func (self *StringBuffer) WriteEnding() {
	self.WriteString(";\r\n")
}

func NewZeusConverter() *ZeusConverter {
	c := ZeusConverter{}
	return &c
}

func (self *ZeusConverter) Parse(filePath string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	parser, err := ParseZeusFile(filePath)
	if err != nil {
		return err
	}

	p := ZeusProtoParser{
		lastError:   nil,
		pbBuffer:    StringBuffer{strings.Builder{}},
		httpContext: &HttpContext{PackageName: "proto"},
	}

	antlr.ParseTreeWalkerDefault.Walk(&p, parser.Proto())

	self.parser = &p
	return nil
}

func (self *ZeusConverter) PBContent() (string, error) {
	return self.parser.pbBuffer.String(), nil
}

func (self *ZeusConverter) HttpContent() (string, error) {
	return self.parser.httpContext.Generate()
}

func ParseZeusFile(path string) (*parser.ZeusProtoParser, error) {
	input, err := antlr.NewFileStream(path)
	if err != nil {
		return nil, err
	}
	lexer := parser.NewZeusProtoLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewZeusProtoParser(stream)
	return p, nil
}

func formatPb(reader *strings.Reader) (string, error) {

	pb, err := proto.NewParser(reader).Parse()
	if err != nil {
		return "", err
	}
	buff := new(bytes.Buffer)
	f := protofmt.NewFormatter(buff, "	")

	f.Format(pb)
	return buff.String(), nil
}

func defaultPbImport() string {
	return `import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "google/protobuf/empty.proto";
import "google/protobuf/wrappers.proto";
import "google/api/annotations.proto";
`
}

func defaultOption() string {
	return `option go_package = "proto";
option (gogoproto.goproto_getters_all) = false;
option (gogoproto.unmarshaler_all) = false;
option (gogoproto.marshaler_all) = false;
option (gogoproto.sizer_all) = false;
`
}
