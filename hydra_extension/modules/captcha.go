package modules

import (
	"github.com/mojocn/base64Captcha"
	"yym/hydra_extension/hydra"
)

func init() {
	hydra.RegisterType("CaptchaGenerator", newCaptchaGenerator)
}

type captchaGenerator struct {
	config interface{}
}

func newCaptchaGenerator() (*captchaGenerator, error) {
	return &captchaGenerator{config: nil}, nil
}
func (self *captchaGenerator) Digit(config base64Captcha.ConfigDigit) {
	self.config = config
}

func (self *captchaGenerator) Character(config base64Captcha.ConfigCharacter) {
	self.config = config
}

func (self *captchaGenerator) Generate(args ...bool) (string, []byte) {
	base64 := false
	if len(args) > 0 {
		base64 = args[0]
	}
	id, data := base64Captcha.GenerateCaptcha("", self.config)
	if base64 {
		return id, []byte(base64Captcha.CaptchaWriteToBase64Encoding(data))
	} else {
		return id, data.BinaryEncoding()
	}
}

func (self *captchaGenerator) Verify(id, value string) bool {
	return base64Captcha.VerifyCaptcha(id, value)
}
