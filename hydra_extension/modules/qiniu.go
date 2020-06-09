package modules

import (
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/sms"
	"yym/hydra_extension/hydra"
)

type authInfo struct {
	Key    string
	Secret string
}

type smsRequest struct {
	SignatureID string
	TemplateID  string
	Mobiles     []string
	Parameters  map[string]string
}

func init() {
	hydra.RegisterType("QiNiuSms", newQiNiuSms)
}

func newQiNiuSms() (*qiniuSms, error) {
	return &qiniuSms{}, nil
}

type qiniuSms struct {
	m *sms.Manager
}

func (self *qiniuSms) SetAuth(au *authInfo) {
	qAuth := auth.New(au.Key, au.Secret)
	self.m = sms.NewManager(qAuth)
}

func (self *qiniuSms) Send(r smsRequest) (sms.MessagesResponse, error) {
	pars := map[string]interface{}{}
	for key, value := range r.Parameters {
		pars[key] = value
	}
	s := sms.MessagesRequest{
		SignatureID: r.SignatureID,
		TemplateID:  r.TemplateID,
		Mobiles:     r.Mobiles,
		Parameters:  pars,
	}

	return self.m.SendMessage(s)
}
