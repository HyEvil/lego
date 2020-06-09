// code from https://code.snowplus.cn/snowplus/ginc/blob/master/util/proto.go

package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/emicklei/proto"
)

// Service represents a service definition is a proto file.
type Service struct {
	Name string
	File string
	APIs []*API
}

// API represents a rpc definition in a proto file, corresponding to a gin controller.
type API struct {
	Name   string
	Method string
	URL    string
	// Target file is the go file where it should be generated.
	TargetFile string
	Request    *Message
	Response   *Message
}

// Message represents a message definition in a proto file.
type Message struct {
	Name string
	// File is the proto file where it's defined.
	File   string
	Fields []*Field
}

// Field represents a field definition in a proto message.
// We only need its inline comment to modify go tags so we don't parse the comment.
type Field struct {
	Name          string
	InlineComment string
}

// ParseProtos parses the proto files and returns a list of Service instances and Message instances.
func ParseProtos(protos []string) ([]*Service, []*Message) {
	msgs := make(map[string]*Message)
	srvs := make([]*Service, 0)
	for _, p := range protos {
		parseProto(p, msgs, &srvs)
	}
	// fill Request and Response of APIs
	for _, srv := range srvs {
		for _, api := range srv.APIs {
			req := msgs[api.Request.Name]
			if req == nil {
				panic(fmt.Sprintf("not found definition `%v` of %v's request", api.Request.Name, api.Name))
			}
			api.Request = req

			res := msgs[api.Response.Name]
			if res == nil {
				panic(fmt.Sprintf("not found definition `%v` of %v's response", api.Response.Name, api.Name))
			}
			api.Response = res
		}

	}

	var messages []*Message
	for _, m := range msgs {
		messages = append(messages, m)
	}

	return srvs, messages
}

func parseProto(file string, msgs map[string]*Message, srvs *[]*Service) {
	reader, _ := os.Open(file)
	defer reader.Close()

	p := proto.NewParser(reader)
	d, err := p.Parse()
	if err != nil {
		panic(err)
	}

	msgPrefix := ""
	if strings.Contains(file, "common.proto") {
		msgPrefix = "common."
	}

	for _, e := range d.Elements {
		switch v := e.(type) {
		case *proto.Service:
			*srvs = append(*srvs, parseService(v, file))
		case *proto.Message:
			msgs[msgPrefix+v.Name] = parseMessage(v, file)
		}
	}
}

func parseService(ps *proto.Service, file string) *Service {
	s := &Service{
		Name: ps.Name,
		File: file,
	}
	var targetFile = "controller.go"
	for _, e := range ps.Elements {
		switch v := e.(type) {
		case *proto.Comment:
			if len(v.Lines) > 0 {
				//windows下行尾有/r,HasSuffix会返回false
				firstLine := strings.TrimSpace(v.Lines[0])
				if strings.HasSuffix(firstLine, ".go") {
					targetFile = firstLine
				}
			}
		case *proto.RPC:
			if v.Comment != nil && len(v.Comment.Lines) > 0 {
				firstLine := strings.TrimSpace(v.Comment.Lines[0])
				if strings.HasSuffix(firstLine, ".go") {
					targetFile = firstLine
				}
			}
			s.APIs = append(s.APIs, parseRPC(v, targetFile))
		}
	}
	return s
}

func parseRPC(pr *proto.RPC, file string) *API {
	a := &API{
		Name:       pr.Name,
		TargetFile: file,
		Request:    &Message{Name: pr.RequestType},
		Response:   &Message{Name: pr.ReturnsType},
	}
	for _, e := range pr.Elements {
		switch v := e.(type) {
		case *proto.Option:
			for k, vv := range v.Constant.Map {
				if vv.IsString && (k == "get" || k == "post" || k == "put" || k == "patch" || k == "delete") {
					a.Method = k
					a.URL = vv.Source
				}
			}
		}
	}
	return a
}

func parseMessage(pm *proto.Message, file string) *Message {
	m := &Message{
		Name: pm.Name,
		File: file,
	}
	for _, e := range pm.Elements {
		switch v := e.(type) {
		case *proto.NormalField:
			field := &Field{Name: v.Name}
			if v.InlineComment != nil && len(v.InlineComment.Lines) > 0 {
				field.InlineComment = strings.TrimSpace(v.InlineComment.Lines[0])
			}
			m.Fields = append(m.Fields, field)
		}
	}
	return m
}
