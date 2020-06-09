package parser

import (
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

type Kind string

const (
	UnKnown Kind = "UnKnown"
	Builtin Kind = "Builtin"
	Struct  Kind = "Struct"
	Map     Kind = "Map"
	Slice   Kind = "Slice"
	Pointer Kind = "Pointer"
	Named   Kind = "Named"
)

type Type interface {
	Kind() Kind
	IsKind(k Kind) bool
	Underlying() *TypeWrapper
	Typ() types.Type
	String() string
	// 不建议用。。。
	Name(withPackage bool) string
	AssignableTo(t Type) bool
	ConvertibleTo(t Type) bool
	GetNamedType() *types.Named
	GetPointerType() *types.Pointer
	GetSliceType() *types.Slice
	AsStructWrapper() *StructWrapper
	AsNamedWrapper() *NamedWrapper
	AsBasicWrapper() *BasicWrapper
	AsPointerWrapper() *PointerWrapper
	AsSliceWrapper() *SliceWrapper
}

type TypeWrapper struct {
	t types.Type
}

type StructWrapper struct {
	*TypeWrapper
	s        *types.Struct
	fieldMap map[string]*TypeWrapper
}

type PointerWrapper struct {
	*TypeWrapper
	p *types.Pointer
}

type BasicWrapper struct {
	*TypeWrapper
	b *types.Basic
}

type NamedWrapper struct {
	*TypeWrapper
	n *types.Named
}

type SliceWrapper struct {
	*TypeWrapper
	s *types.Slice
}

type Packages []*packages.Package

func NewTypeWrapper(t types.Type) *TypeWrapper {
	return &TypeWrapper{t: t}
}

func (self *Packages) Lookup(path string) (*TypeWrapper, error) {
	for _, pkg := range *self {
		if pkg.Errors != nil {
			continue
		}
		if pkg.Types == nil {
			continue
		}
		obj := pkg.Types.Scope().Lookup(path)
		if obj != nil {
			return &TypeWrapper{t: obj.Type()}, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

type StructFiledFilter func(string) bool

func MakeFieldExportedFilter() StructFiledFilter {
	return func(name string) bool {
		return token.IsExported(name)
	}
}

func MakeFieldNotContainsFilter(key string) StructFiledFilter {
	return func(name string) bool {
		return strings.Index(name, key) == -1
	}
}

func (self *TypeWrapper) Kind() Kind {
	switch self.t.(type) {
	case *types.Basic:
		return Builtin
	case *types.Pointer:
		return Pointer
	case *types.Named:
		return Named
	case *types.Map:
		return Map
	case *types.Struct:
		return Struct
	case *types.Slice:
		return Slice
	default:
		return UnKnown
	}
}

func (self *TypeWrapper) IsKind(k Kind) bool {
	return self.Kind() == k
}

func (self *TypeWrapper) Underlying() *TypeWrapper {
	return &TypeWrapper{t: self.t.Underlying()}
}

func (self *TypeWrapper) Typ() types.Type {
	return self.t
}

func (self *TypeWrapper) String() string {
	return self.t.String()
}

func (self *TypeWrapper) AssignableTo(t Type) bool {
	if self.Kind() == Named && t.Kind() == Named {
		if self.String() == t.String() {
			return true
		}
	}
	return types.AssignableTo(self.t, t.Typ())
}

func (self *TypeWrapper) ConvertibleTo(t Type) bool {
	return types.ConvertibleTo(self.t, t.Typ())
}

func (self *TypeWrapper) GetNamedType() *types.Named {
	return self.t.(*types.Named)
}

func (self *TypeWrapper) GetPointerType() *types.Pointer {
	return self.t.(*types.Pointer)
}

func (self *TypeWrapper) GetSliceType() *types.Slice {
	return self.t.(*types.Slice)
}

func (self *TypeWrapper) AsBasicWrapper() *BasicWrapper {
	return &BasicWrapper{TypeWrapper: &TypeWrapper{t: self.t}, b: self.t.(*types.Basic)}
}

func (self *TypeWrapper) AsStructWrapper() *StructWrapper {
	return &StructWrapper{TypeWrapper: &TypeWrapper{t: self.t}, s: self.t.(*types.Struct)}
}

func (self *TypeWrapper) AsNamedWrapper() *NamedWrapper {
	return &NamedWrapper{TypeWrapper: &TypeWrapper{t: self.t}, n: self.t.(*types.Named)}
}

func (self *TypeWrapper) AsSliceWrapper() *SliceWrapper {
	return &SliceWrapper{TypeWrapper: &TypeWrapper{t: self.t}, s: self.t.(*types.Slice)}
}

func (self *TypeWrapper) Name(withPackage bool) string {
	match := regexp.MustCompile("([\\[\\*\\]]*)(?:[\\w]+/)*(?:([\\w]+)\\.)?(\\w+)").FindStringSubmatch(self.String())
	if len(match) == 0 {
		return ""
	}
	if withPackage {
		return match[1] + match[2] + "." + match[3]
	}
	return match[1] + match[3]
}

func (self *TypeWrapper) AsPointerWrapper() *PointerWrapper {
	return &PointerWrapper{TypeWrapper: &TypeWrapper{t: self.t}, p: self.t.(*types.Pointer)}
}

func (self *StructWrapper) Type() *types.Struct {
	return self.s
}

func (self *StructWrapper) NumFields() int {
	return self.s.NumFields()
}

func (self *StructWrapper) Field(i int) (string, *TypeWrapper) {
	return self.s.Field(i).Name(), &TypeWrapper{t: self.s.Field(i).Type()}
}

func (self *StructWrapper) FieldByName(name string) *TypeWrapper {
	m := self.initFieldMap()
	return m[name]
}

func (self *StructWrapper) initFieldMap() map[string]*TypeWrapper {
	if self.fieldMap != nil {
		return self.fieldMap
	}
	self.fieldMap = map[string]*TypeWrapper{}
	for i := 0; i < self.s.NumFields(); i++ {
		n, f := self.Field(i)
		self.fieldMap[n] = f
	}
	return self.fieldMap
}

func (self *SliceWrapper) Elem() *TypeWrapper {
	return &TypeWrapper{t: self.s.Elem()}
}

// 根据过滤器生成一个新的struct type
func (self *StructWrapper) NewStructWithFilter(filters ...StructFiledFilter) *StructWrapper {
	vars := []*types.Var{}
	tags := []string{}

	for i := 0; i < self.s.NumFields(); i++ {
		isRemove := false
		field := self.s.Field(i)

		for _, filter := range filters {
			if filter != nil && !filter(field.Name()) {
				isRemove = true
				break
			}
		}
		if isRemove {
			continue
		}
		vars = append(vars, field)
		tags = append(tags, self.s.Tag(i))
	}
	s := types.NewStruct(vars, tags)
	return &StructWrapper{TypeWrapper: &TypeWrapper{t: s,}, s: s}
}

func (self *StructWrapper) FuzzyField(name string) (string, *TypeWrapper) {
	fieldNameList := []string{}
	for i := 0; i < self.s.NumFields(); i++ {
		fieldNameList = append(fieldNameList, self.s.Field(i).Name())
	}
	matches := fuzzy.Find(name, fieldNameList)
	if matches.Len() == 0 {
		return "", nil
	}

	return matches[0].Str, self.FieldByName(matches[0].Str)

}

func (self *PointerWrapper) Type() *types.Pointer {
	return self.p
}

func (self *PointerWrapper) Elem() *TypeWrapper {
	return &TypeWrapper{t: self.p.Elem()}
}

func (self *BasicWrapper) Type() *types.Basic {
	return self.b
}

func (self *BasicWrapper) IsNumeric() bool {
	return self.b.Info()&types.IsNumeric != 0
}

func (self *BasicWrapper) IsString() bool {
	return self.b.Info()&types.IsString != 0
}

func (self *BasicWrapper) IsConstType() bool {
	return self.b.Info()&types.IsConstType != 0
}

func (self *NamedWrapper) Type() *types.Named {
	return self.n
}

/*
func parseLocal(path string) (pkgs map[string]*ast.Package, first error) {
	fset := token.NewFileSet()
	c := types.Config{
		IgnoreFuncBodies: true,
		// Note that importAdapter can call b.importPackage which calls this
		// method. So there can't be cycles in the import graph.
		Importer: importer.Default(),
		Error: func(err error) {
		},
	}
	return parser.ParseDir(fset, path, nil, 0)
}

func parseTypeFromLocal(packagePath, name string) (*TypeWrapper, error) {
	pkgs, err := parseLocal(packagePath)
	if err != nil {
		return nil, err
	}

}
*/
func ParsePackages(path, dir string) (Packages, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportsFile |
			packages.NeedTypes |
			//packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
		Dir: dir,
	}, path)
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func ParseTypeFromPackage(path, dir, name string) (*TypeWrapper, error) {
	pkgs, err := ParsePackages(path, dir)
	if err != nil {
		return nil, err
	}

	var obj types.Object
	for _, pkg := range pkgs {
		if pkg.Types == nil || pkg.Types.Scope() == nil {
			continue
		}
		obj = pkg.Types.Scope().Lookup(name)
		if obj != nil {
			break
		}
	}
	if obj == nil {
		return nil, fmt.Errorf("not found")
	}
	t := &TypeWrapper{t: obj.Type()}
	return t, nil
}

func ParseTypeByFileLineOffset(path, dir, filePath string, line, offset int) (*TypeWrapper, error) {
	logrus.Debug("ParseTypeByFileLineOffset filepath %s", filePath)
	pkgs, err := ParsePackages(path, dir)
	if err != nil {
		return nil, err
	}

	var goFile *token.File
	for _, pkg := range pkgs {
		if pkg.Errors != nil {
			for _, e := range pkg.Errors {
				logrus.WithError(e).Errorf("parse pkg %s error", pkg.Name)
			}
		}
		pkg.Fset.Iterate(func(file *token.File) bool {
			if filepath.ToSlash(file.Name()) == filepath.ToSlash(filePath) {
				goFile = file
				return false
			}
			return true
		})

	}
	if goFile == nil {
		return nil, fmt.Errorf("file not found:%s", filePath)
	}

	pos := goFile.LineStart(line) + token.Pos(offset)

	for _, pkg := range pkgs {

		for expr, t := range pkg.TypesInfo.Types {
			if expr.Pos() == pos {
				return &TypeWrapper{t: t.Type}, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

func ParseTypeByFileOffset(path, dir, filePath string, offset int) (*TypeWrapper, error) {
	pkgs, err := ParsePackages(path, dir)
	if err != nil {
		return nil, err
	}

	var goFile *token.File
	for _, pkg := range pkgs {
		pkg.Fset.Iterate(func(file *token.File) bool {
			name := filepath.Base(file.Name())
			if name == filePath {
				goFile = file
				return false
			}
			return true
		})

	}
	if goFile == nil {
		return nil, fmt.Errorf("file not found")
	}
	pos := goFile.Pos(offset)

	for _, pkg := range pkgs {

		for expr, t := range pkg.TypesInfo.Types {
			if expr.Pos() == pos {
				return &TypeWrapper{t: t.Type}, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}
