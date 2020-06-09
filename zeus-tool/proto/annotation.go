package proto

import (
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"go.yym.plus/zeus-tool/proto/parser"
	"go.yym.plus/zeus-tool/utils"
)

type AnnotationType int

const (
	AnnotationTypeService AnnotationType = iota + 1
	AnnotationTypeRpc
	AnnotationTypeMessage
	AnnotationTypeField
)

type AnnotationProcessor func(converter *ZeusProtoParser, ctx antlr.ParserRuleContext)

var (
	messageAnnotationProcessors = map[AnnotationType]map[string]AnnotationProcessor{}
)

func init() {
	registerAnnotationProcessor(AnnotationTypeField, "tag", onTagAnnotation)
	registerAnnotationProcessor(AnnotationTypeService, "resource", onResourceAnnotation)
	registerAnnotationProcessor(AnnotationTypeRpc, "middleware", onRpcMiddlewareAnnotation)
	registerAnnotationProcessor(AnnotationTypeService, "middleware", onServiceMiddlewareAnnotation)
	registerAnnotationProcessor(AnnotationTypeMessage, "inheritable", onInheritableAnnotation)
	registerAnnotationProcessor(AnnotationTypeRpc, "get", onRequestAnnotation)
	registerAnnotationProcessor(AnnotationTypeRpc, "post", onRequestAnnotation)
	registerAnnotationProcessor(AnnotationTypeRpc, "put", onRequestAnnotation)
	registerAnnotationProcessor(AnnotationTypeRpc, "delete", onRequestAnnotation)
}

func registerAnnotationProcessor(t AnnotationType, name string, processor func(converter *ZeusProtoParser, ctx antlr.ParserRuleContext)) {
	if messageAnnotationProcessors[t] == nil {
		messageAnnotationProcessors[t] = map[string]AnnotationProcessor{}
	}
	messageAnnotationProcessors[t][name] = processor
}

func onTagAnnotation(converter *ZeusProtoParser, ctx antlr.ParserRuleContext) {
	switch ann := ctx.(*parser.AnnotationContext).AnnotationWithNamedArg().(type) {
	case *parser.AnnotationWithNamedArgContext:
		for _, arg := range ann.AllAnnotationNamedArg() {
			key := arg.(*parser.AnnotationNamedArgContext).AnnotationNamedArgKey().GetText()
			value := arg.(*parser.AnnotationNamedArgContext).AnnotationNamedArgValue().GetText()
			switch key {
			case "default":
				converter.ctx.curField.defaultTag = value
			case "validator":
				converter.ctx.curField.validatorTag = value
			case "json":
				converter.ctx.curField.jsonTag = value
			case "customtype":
				converter.ctx.curField.customTag = value
			case "nullable":
				converter.ctx.curField.nullable = value
			default:
				panic(utils.WithPosError(ctx, "tag annotation arg %s not support", key))
			}
		}
	}
}

func onResourceAnnotation(converter *ZeusProtoParser, ctx antlr.ParserRuleContext) {
	switch ann := ctx.(*parser.AnnotationContext).AnnotationWithArg().(type) {
	case *parser.AnnotationWithArgContext:
		var err error
		converter.ctx.curService.resource, err = strconv.Unquote(ann.AllAnnotationArg()[0].GetText())
		if err != nil {
			panic(utils.WithPosError(ctx, "%v", err))
		}
	default:
		panic(utils.WithPosError(ctx, "resource annotation must have arg"))
	}
	//converter.ctx.curService.resource = ctx.
}

func onRpcMiddlewareAnnotation(converter *ZeusProtoParser, ctx antlr.ParserRuleContext) {
	switch ann := ctx.(*parser.AnnotationContext).AnnotationWithArg().(type) {
	case *parser.AnnotationWithArgContext:
		converter.ctx.curRpc.middlewareList = append(converter.ctx.curRpc.middlewareList, ann.AllAnnotationArg()[0].GetText())
	default:
		panic(utils.WithPosError(ctx, "resource annotation must have arg"))
	}
}

func onServiceMiddlewareAnnotation(converter *ZeusProtoParser, ctx antlr.ParserRuleContext) {
	switch ann := ctx.(*parser.AnnotationContext).AnnotationWithArg().(type) {
	case *parser.AnnotationWithArgContext:
		converter.ctx.curService.middlewareList = append(converter.ctx.curService.middlewareList, ann.AllAnnotationArg()[0].GetText())
	default:
		panic(utils.WithPosError(ctx, "resource annotation must have arg"))
	}
}

func onRequestAnnotation(converter *ZeusProtoParser, ctx antlr.ParserRuleContext) {
	switch ann := ctx.(*parser.AnnotationContext).AnnotationWithArg().(type) {
	case *parser.AnnotationWithArgContext:
		method := ann.AnnotationName().GetText()
		converter.ctx.curRpc.httpMethod = method
		if len(ann.AllAnnotationArg()) != 1 {
			panic(utils.WithPosError(ctx, "%s annotation arg count error", method))
		}
		var err error
		converter.ctx.curRpc.httpPath, err = strconv.Unquote(ann.AllAnnotationArg()[0].GetText())
		if err != nil {
			panic(utils.WithPosError(ctx, "%v", err))
		}
	}
}

func onInheritableAnnotation(converter *ZeusProtoParser, ctx antlr.ParserRuleContext) {

}
