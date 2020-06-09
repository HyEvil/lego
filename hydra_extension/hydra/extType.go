package hydra

import (
	"github.com/spf13/cast"
	"yym/hydra_extension/3rd/codec"
	"time"
)

type Duration time.Duration

func (self *Duration) CodecEncodeSelf(en *codec.Encoder) {
	str, err := cast.ToStringE(time.Duration(*self))
	if err != nil {
		panic(err)
	}
	en.Encode(str)
}

func (self *Duration) CodecDecodeSelf(de *codec.Decoder) {
	var s string
	de.MustDecode(&s)
	d, err := cast.ToDurationE(s)
	if err != nil {
		panic(err)
	}
	*self = Duration(d)
}

func (self *Duration) Value() time.Duration {
	return time.Duration(*self)
}
