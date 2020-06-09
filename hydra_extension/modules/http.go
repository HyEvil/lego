package modules

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/spf13/cast"
	"net/http"
	"net/url"
	"time"
	"yym/hydra_extension/hydra"
)

func init() {
	hydra.RegisterType("HttpClient", NewHttpClient)
}

type HttpRequest struct {
	Url    string                  `codec:"url"`
	Header *map[string]string      `codec:"header"`
	Param  *map[string]interface{} `codec:"param"`
	Data   *[]byte                 `codec:"data"`
}

type HttpResponse struct {
	Code   int               `codec:"code"`
	Header map[string]string `codec:"header"`
	Data   []byte            `codec:"data"`
	Cost   float64           `codec:"cost"`
	Debug  string            `codec:"debug"`
}

func (self *HttpRequest) toArgs() []interface{} {
	var args []interface{}
	if self.Header != nil {
		args = append(args, req.Header(*self.Header))
	}
	if self.Param != nil {
		p := req.Param{}
		for key, value := range *self.Param {
			p[key] = cast.ToString(value)
		}
		args = append(args, p)
	}
	if self.Data != nil {
		args = append(args, *self.Data)
	}
	return args
}

type HttpClient struct {
	r     *req.Req
	debug bool
}

func NewHttpClient() (*HttpClient, error) {
	r := req.New()
	r.SetFlags(req.LstdFlags | req.Lcost)
	return &HttpClient{
		r:     r,
		debug: false,
	}, nil
}

func (self *HttpClient) Get(request HttpRequest) (*HttpResponse, error) {
	return self.toResp(self.r.Get(request.Url, request.toArgs()...))
}

func (self *HttpClient) Post(request HttpRequest) (*HttpResponse, error) {
	return self.toResp(self.r.Post(request.Url, request.toArgs()...))
}

func (self *HttpClient) Delete(request HttpRequest) (*HttpResponse, error) {
	return self.toResp(self.r.Delete(request.Url, request.toArgs()...))
}

func (self *HttpClient) Head(request HttpRequest) (*HttpResponse, error) {
	return self.toResp(self.r.Head(request.Url, request.toArgs()...))
}

func (self *HttpClient) Put(request HttpRequest) (*HttpResponse, error) {
	return self.toResp(self.r.Put(request.Url, request.toArgs()...))
}

func (self *HttpClient) EnableCookie(b bool) {
	self.r.EnableCookie(b)
}

func (self *HttpClient) EnableDebug(b bool) {
	self.debug = b
}

func (self *HttpClient) SetTimeout(f hydra.Duration) {
	self.r.SetTimeout(f.Value())
}

func (self *HttpClient) SetProxy(p string) error {
	if p == "" {
		return self.r.SetProxy(func(request *http.Request) (i *url.URL, e error) {
			return nil, nil
		})
	} else {
		return self.r.SetProxy(func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(p)
		})
	}
}

func (self *HttpClient) toResp(r *req.Resp, err error) (*HttpResponse, error) {
	if r == nil {
		return nil, err
	}

	resp := &HttpResponse{
		Code:   r.Response().StatusCode,
		Data:   r.Bytes(),
		Header: map[string]string{},
		Cost:   float64(r.Cost()) / float64(time.Second),
	}
	for key, value := range r.Response().Header {
		if len(value) > 0 {
			resp.Header[key] = value[0]
		}
	}
	if self.debug {
		resp.Debug = fmt.Sprintf("%+v", r)
	}
	return resp, err
}
