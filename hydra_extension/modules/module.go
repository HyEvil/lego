package modules

import "yym/hydra_extension/hydra"

type module struct {
}

func (self *module) Name() string {
	return "builtin"
}

func (self *module) Version() string {
	return "0.0.2"
}

func init() {
	hydra.RegisterModule("module", &module{})
}
