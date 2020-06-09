package generator

import (
	"fmt"
	"go/types"
	"regexp"
	"strings"
	"yym/snowpluslib/parser"
	"yym/snowpluslib/utils"
)

var (
	errCanNotConvert = fmt.Errorf("can not convert")
	errNotExported   = fmt.Errorf("not exported")
	errNotSupport    = fmt.Errorf("can support")
)

type WrapType string

const (
	WrapOther        WrapType = "other"
	WrapCastCall     WrapType = "funcCall"
	WrapConvert      WrapType = "convert"
	WrapBasicConvert WrapType = "basicConvert"
	WrapLambda       WrapType = "lambda"
)

type TypeTransformer struct {
	handlers      map[parser.Kind]func(source, target parser.Type, varName string) (string, error)
	specConvert   []func(source parser.Type, varName string) (parser.Type, string, error)
	isSamePackage bool
	stack         [][]parser.Kind
}

// todo 可以外边嵌套一层，避免每次init
func NewTypeTransformer(samePackage bool) *TypeTransformer {
	t := &TypeTransformer{handlers: map[parser.Kind]func(source parser.Type, target parser.Type, varName string) (string, error){},
		isSamePackage: samePackage}

	t.Init()
	return t
}

func GenerateStructTransform(isSamePkg bool, sourceStruct, targetStruct *parser.StructWrapper, sourceVarName string, sourceFilter, targetFilter parser.StructFiledFilter) (string, error) {
	if !sourceStruct.IsKind(parser.Struct) || !targetStruct.IsKind(parser.Struct) {
		return "", fmt.Errorf("source and target must be struct type")
	}

	buffer := strings.Builder{}
	writeAssign := func(left, right string, needComment bool) {

		if needComment {
			buffer.WriteString("//")
			right = strings.Replace(right, utils.LineBreak, utils.LineBreak+"//", -1)
		}

		buffer.WriteString(left + ":")
		buffer.WriteString(right + ",")
		buffer.WriteString(utils.LineBreak)
	}
	source := sourceStruct.NewStructWithFilter(parser.MakeFieldExportedFilter(), sourceFilter)
	target := targetStruct.NewStructWithFilter(parser.MakeFieldExportedFilter(), targetFilter)

	typeTransformer := NewTypeTransformer(isSamePkg)
	//buffer.WriteString(targetStruct.GetNamedType().Obj().Name() + "{" + utils.LineBreak)

	for i := 0; i < target.NumFields(); i++ {
		targetFieldName, targetField := target.Field(i)

		sourceField := source.FieldByName(targetFieldName)

		if sourceField != nil {
			rightVarName := sourceVarName + "." + targetFieldName
			rightExpr, needComment := typeTransformer.transform(sourceField, targetField, rightVarName, true)
			writeAssign(targetFieldName, rightExpr, needComment)
			continue
		}

		sourceFieldName, sourceField := source.FuzzyField(targetFieldName)
		if sourceField != nil {
			//todo 匹配值匹配度检查
			rightVarName := sourceVarName + "." + sourceFieldName
			rightExpr, _ := typeTransformer.transform(sourceField, targetField, rightVarName, true)
			writeAssign(targetFieldName, rightExpr, true)
			continue
		}

		writeAssign(targetFieldName, getTypeDefaultValue(targetField), true)
	}
	//buffer.WriteString("}")
	return buffer.String(), nil
}

func getTypeDefaultValue(v parser.Type) string {
	underlying := v.Underlying()
	if underlying.IsKind(parser.Builtin) {
		if underlying.AsBasicWrapper().IsNumeric() {
			return "0"
		} else if underlying.AsBasicWrapper().IsString() {
			return "\"\""
		} else {
			return "nil"
		}
	} else {
		return "nil"
	}
}

func (self *TypeTransformer) Init() {
	//todo struct to struct得支持一哈。。。
	self.register("*", parser.Named, self.AnyToNamed)
	self.register(parser.Named, "*", self.NamedToAny)
	self.register(parser.Builtin, parser.Builtin, self.BuiltinToBuiltin)
	self.register("*", parser.Pointer, self.AnyToPointer)
	self.register(parser.Pointer, "*", self.PointerToAny)
	self.register(parser.Slice, parser.Slice, self.SliceToSlice)

	self.registerSpecConvert(self.ConvertOriginal)
	self.registerSpecConvert(self.ConvertFromTimeSpec)
	self.registerSpecConvert(self.ConvertToTimeSpec)
}

func (self *TypeTransformer) register(sourceType, targetType parser.Kind, f func(source, target parser.Type, varName string) (string, error)) {
	self.handlers[sourceType+targetType] = f
}

func (self *TypeTransformer) registerSpecConvert(f func(source parser.Type, varName string) (parser.Type, string, error)) {
	self.specConvert = append(self.specConvert, f)
}

func (self *TypeTransformer) transform(source, target parser.Type, varName string, match bool) (string, bool) {
	if source.AssignableTo(target) {
		return varName, false
	}
	expr, err := self.generateRightAssignment(source, target, varName)
	if err != nil {
		return getTypeDefaultValue(target), true
	}
	return expr, true
}

func (self *TypeTransformer) generateRightAssignment(source, target parser.Type, varName string) (string, error) {

	for _, convert := range self.specConvert {
		newSource, newSourceVar, err := convert(source, varName)
		if err != nil {
			continue
		}
		ret, err := self.generateRightAssignmentEx(newSource, target, newSourceVar)
		if err == nil {
			return ret, nil
		}
	}

	return "", errNotSupport
}

func (self *TypeTransformer) generateRightAssignmentEx(source, target parser.Type, varName string) (string, error) {
	if source.AssignableTo(target) {
		return varName, nil
	}
	handler := self.handlers[source.Kind()+target.Kind()]
	if handler == nil {
		handler = self.handlers[source.Kind()+parser.Kind("*")]
	}
	if handler == nil {
		handler = self.handlers[parser.Kind("*")+target.Kind()]
	}
	if handler == nil {
		return "", fmt.Errorf("no handler")
	}

	self.stack = append(self.stack, []parser.Kind{source.Kind(), target.Kind()})
	ret, err := handler(source, target, varName)
	self.stack = self.stack[:len(self.stack)-1]
	return ret, err
}

func (self *TypeTransformer) NamedToNamed(source, target parser.Type, varName string) (string, error) {
	if target.String() == source.String() {
		return varName, nil
	}
	return "", errNotSupport
}

func (self *TypeTransformer) AnyToNamed(source, target parser.Type, varName string) (string, error) {
	if !target.AsNamedWrapper().Type().Obj().Exported() {
		return "", errNotExported
	}
	expr, err := self.generateRightAssignment(source, target.Underlying(), varName)
	if err != nil {
		return "", err
	}
	ret := ""
	if !self.isSamePackage {
		ret = target.AsNamedWrapper().Type().Obj().Pkg().Name() + "."
	}
	if t, _ := self.wrapType(expr); t == WrapConvert {
		index := strings.Index(expr, "(")
		expr = expr[index+1:]
		expr = expr[:len(expr)-1]
		ret = ret + target.AsNamedWrapper().Type().Obj().Name() + "(" + expr + ")"
	} else {
		ret = ret + target.AsNamedWrapper().Type().Obj().Name() + "(" + expr + ")"
	}

	return ret, nil
}

func (self *TypeTransformer) NamedToAny(source, target parser.Type, varName string) (string, error) {
	expr, err := self.generateRightAssignment(source.Underlying(), target, source.Underlying().String()+"("+varName+")")
	return expr, err
}

func (self *TypeTransformer) AnyToPointer(source, target parser.Type, varName string) (string, error) {
	expr, err := self.generateRightAssignment(source, target.AsPointerWrapper().Elem(), varName)
	if err != nil {
		return "", err
	}
	ret := ""

	t, sub := self.wrapType(expr)

	if t == WrapCastCall {
		ret = fmt.Sprintf(`func()*%s{
	v:=%s
	return &v}()`, sub[0], expr)

	} else if t == WrapBasicConvert || t == WrapConvert {
		ret = fmt.Sprintf(`func()*%s{
	v:=%s
	return &v}()`, sub[0], expr)

	} else if t == WrapLambda {
		ret = strings.Replace(expr, "func()", "func()*", 1)
		retExprPos := strings.LastIndex(ret, "v}()")
		ret = ret[:retExprPos] + "&" + ret[retExprPos:]
	} else {

		ret = ret + "&" + expr
	}

	return ret, nil
}

func (self *TypeTransformer) PointerToAny(source, target parser.Type, varName string) (string, error) {
	expr, err := self.generateRightAssignment(source.AsPointerWrapper().Elem(), target, "*"+varName)
	if err != nil {
		return "", err
	}

	return expr, nil
}

func (self *TypeTransformer) BuiltinToBuiltin(source, target parser.Type, varName string) (string, error) {
	sourceBasic := source.AsBasicWrapper()
	targetBasic := target.AsBasicWrapper()
	if sourceBasic.IsString() && targetBasic.IsString() {
		return varName, nil
	} else if sourceBasic.IsNumeric() && targetBasic.IsNumeric() {
		return targetBasic.Type().Name() + "(" + varName + ")", nil
	}
	return "cast.To" + strings.Title(targetBasic.Type().Name()) + "(" + varName + ")", nil
}

func (self *TypeTransformer) StructToStruct(source, target parser.Type, varName string) (string, error) {
	return "", errNotSupport
}

func (self *TypeTransformer) SliceToSlice(source, target parser.Type, varName string) (string, error) {
	sourceElem := source.AsSliceWrapper().Elem()
	targetElem := target.AsSliceWrapper().Elem()
	expr, err := self.generateRightAssignment(sourceElem, targetElem, "e")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`func()%s{
s:=%s{}
for _,e:=range %s{
	s= append(s,%s)
}
return s
}()`, target.Name(!self.isSamePackage), target.Name(!self.isSamePackage), varName, expr), nil

}

func (self *TypeTransformer) preStack() []parser.Kind {
	if len(self.stack) < 2 {
		return nil
	}
	return self.stack[len(self.stack)-2]
}

func (self *TypeTransformer) wrapType(expr string) (WrapType, []string) {
	if regexp.MustCompile("^func\\(").MatchString(expr) {
		return WrapLambda, nil
	} else if sub := regexp.MustCompile("To\\w+\\(\\)$").FindStringSubmatch(expr); len(sub) > 0 {
		return WrapCastCall, []string{utils.LowerFirst(sub[1])}
	} else if sub := regexp.MustCompile("^(\\w+)\\(").FindStringSubmatch(expr); len(sub) > 0 {
		if self.isBuiltInType(sub[1]) {
			return WrapBasicConvert, sub[1:]
		}
		return WrapConvert, sub[1:]
	}
	return WrapOther, nil
}

func (self *TypeTransformer) ConvertOriginal(source parser.Type, varName string) (parser.Type, string, error) {
	return source, varName, nil
}

func (self *TypeTransformer) ConvertFromTimeSpec(source parser.Type, varName string) (parser.Type, string, error) {
	if !source.IsKind(parser.Named) {
		return nil, "", errNotSupport
	}
	structType := source.AsNamedWrapper()
	if structType.Name(true) != "time.Time" {
		return nil, "", errNotSupport
	}

	return parser.NewTypeWrapper(types.Universe.Lookup("int64").Type()), varName + ".Unix()", nil
}

func (self *TypeTransformer) ConvertToTimeSpec(source parser.Type, varName string) (parser.Type, string, error) {
	if !source.IsKind(parser.Builtin) {
		return nil, "", errNotSupport
	}
	basic := source.AsBasicWrapper()
	switch basic.Type().Kind() {
	case types.Int32, types.Uint32, types.Uint64:
		return parser.GetCachedType("time.Time", "time", "Time"), fmt.Sprintf("time.Unix(int64(%s),0)", varName), nil
	case types.Int64:
		return parser.GetCachedType("time.Time", "time", "Time"), fmt.Sprintf("time.Unix(%s,0)", varName), nil
	default:
		return nil, "", errNotSupport
	}

}

func (self *TypeTransformer) SpecType(name string) *parser.TypeWrapper {
	return nil

}

// 避免isPrimary转一层。。
func (self *TypeTransformer) isBuiltInType(typ string) bool {
	switch typ {
	case "bool", "byte", "complex128", "complex64":
	case "float32", "float64":
	case "int", "int16", "int32", "int64", "int8":
	case "rune", "string":
	case "uint", "uint16", "uint32", "uint64", "uint8":
	default:
		return false
	}
	return true
}
