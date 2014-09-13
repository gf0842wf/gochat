package codec

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
)

import (
	"net/ttcp/proto"
	"zcodec"
	"zrandom"
)

// prepare decode, 解析长度
func PreDecode(conn *net.TCPConn, header []byte) (data []byte, err error) {
	// header
	// --这个 ReadFull 非常好用, 作用是一直等到读取header大小的字节数为止
	n, err := io.ReadFull(conn, header)
	if err != nil {
		err = errors.New("Error recv header:" + strconv.Itoa(n))
	}

	// data
	length := binary.BigEndian.Uint32(header)
	data = make([]byte, length)
	n, err = io.ReadFull(conn, data)
	if err != nil {
		err = errors.New("Error recv msg:" + strconv.Itoa(n))
	}

	return
}

// after encode, 封装长度
func AftEncode(msg []byte) (data []byte, err error) {
	// msg: 消息, data: 封装后带长度的消息
	length := uint32(len(msg))
	data = make([]byte, length+4)
	binary.BigEndian.PutUint32(data[:4], length)
	data = append(data[:4], msg...)
	return
}

type Coder struct {
	Shaked   bool   // 握手标志
	Encrypt  bool   // 加密标志
	CryptKey uint32 // 加密的key, 由服务器端随机生成发给客户端
}

func NewCoder() *Coder {
	key := zrandom.Randint(0, 2<<31-1) // 随机生成key
	return &Coder{Encrypt: true, CryptKey: uint32(key)}
}

func (dcr *Coder) Decode(msg []byte, obj *proto.Msg) (err error) {
	if dcr.Shaked && dcr.Encrypt {
		zcodec.Crypt(dcr.CryptKey, msg) // 解密
	}
	err = json.Unmarshal(msg, obj)
	return
}

func (dcr *Coder) Encode(obj *proto.Msg) (data []byte, err error) {
	msg, err := obj.Json()
	if dcr.Shaked && dcr.Encrypt {
		zcodec.Crypt(dcr.CryptKey, msg) // 加密
	}
	data, _ = AftEncode(msg)
	return
}
