package proto

import (
	"fmt"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/iancoleman/strcase"
	"github.com/thoas/go-funk"
)

var httpTpl = pongo2.Must(pongo2.FromString(`
{% for service in serviceList %}
func Register{{ service.Name }}HttpServer(rpcServer {{ service.Name }}Server, httpServer *http.Engine) {
	{% for handler in service.Handlers -%}
	httpServer.{{ handler.Method }}("{{ handler.Path }}", {{ mixMiddleware(service.MiddlewareList,handler.MiddlewareList) }} api.HttpWrapper(rpcServer.{{ handler.Name }}))
	{% endfor -%}
}
{% endfor %}
`))

type HttpContext struct {
	PackageName string
	ServiceList []*HttpService
}

type HttpService struct {
	Name           string
	Handlers       []*HttpHandler
	MiddlewareList []string
}

type HttpHandler struct {
	Name           string
	Method         string
	Path           string
	MiddlewareList []string
}

func (self *HttpContext) NewService(name string) {
	s := &HttpService{
		Name: strcase.ToCamel(name),
	}
	self.ServiceList = append(self.ServiceList, s)
}

func (self *HttpContext) NewHandler(name string) {
	s := self.ServiceList[len(self.ServiceList)-1]
	s.Handlers = append(s.Handlers, &HttpHandler{
		Name: strcase.ToCamel(name),
	})
}

func mixMiddleware(mlist ...[]string) *pongo2.Value {
	middlewares := []string{}
	for _, l := range mlist {
		middlewares = append(middlewares, l...)
	}
	m := map[string]bool{}
	for _, s := range middlewares {
		m[s] = true
	}
	ret := ""
	for s, _ := range m {
		ret = ret + s + ","
	}
	return pongo2.AsValue(ret)
}

func (self *HttpContext) Generate() (string, error) {
	for _, service := range self.ServiceList {
		service.Handlers = funk.Filter(service.Handlers, func(h *HttpHandler) bool {
			return h.Method != "" && h.Path != ""
		}).([]*HttpHandler)
	}

	content, err := httpTpl.Execute(pongo2.Context{
		"serviceList":   self.ServiceList,
		"mixMiddleware": mixMiddleware,
	})
	if err != nil {
		return "", err
	}
	buff := strings.Builder{}
	buff.WriteString("// Code generated by zeus-tool. DO NOT EDIT.\n")
	buff.WriteString(fmt.Sprintf("package %s \n", self.PackageName))
	buff.WriteString(content)
	return buff.String(), nil
}

func (self *HttpContext) SetHandlerPath(path string) {
	self.CurrentHandler().Path = path
}

func (self *HttpContext) SetHandlerMethod(method string) {
	self.CurrentHandler().Method = strings.ToUpper(method)
}

func (self *HttpContext) AddServiceMiddleware(m ...string) {
	self.CurrentService().MiddlewareList = append(self.CurrentService().MiddlewareList, m...)
}

func (self *HttpContext) AddHandlerMiddleware(m ...string) {
	self.CurrentHandler().MiddlewareList = append(self.CurrentHandler().MiddlewareList, m...)
}

func (self *HttpContext) CurrentService() *HttpService {
	return self.ServiceList[len(self.ServiceList)-1]
}

func (self *HttpContext) CurrentHandler() *HttpHandler {
	s := self.CurrentService()
	return s.Handlers[len(s.Handlers)-1]
}

func (self *HttpService) AddHandle(method string, path string) {
	self.Handlers = append(self.Handlers, &HttpHandler{
		Method: strings.ToUpper(method),
		Path:   path,
	})
}

func (self *HttpHandler) MiddlewareString() string {
	s := strings.Join(self.MiddlewareList, ",")
	if s != "" {
		return s + ","
	}
	return ""
}
