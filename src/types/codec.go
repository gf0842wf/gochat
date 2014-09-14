package types

import (
	"misc/zcodec"
	"misc/zrandom"
)

type Coder struct {
	Shaked   bool   // 是否握手?
	Encrypt  bool   // 是否需要加密?
	CryptKey uint32 // 加密的key, 由服务器端随机生成发给客户端
}

func NewCoder() *Coder {
	key := zrandom.Randint(1, 2<<31-1) // 随机生成key
	return &Coder{Shaked: false, Encrypt: true, CryptKey: uint32(key)}
}

func (cr *Coder) Decode(data []byte) (err error) {
	if cr.Shaked && cr.Encrypt {
		zcodec.Crypt(cr.CryptKey, data) // 解密
	}

	return
}

func (cr *Coder) Encode(data []byte) (err error) {
	if cr.Shaked && cr.Encrypt {
		zcodec.Crypt(cr.CryptKey, data) // 加密
	}
	return
}
