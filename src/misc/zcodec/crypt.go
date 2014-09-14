package zcodec

// 加密解密, 加密加密为同一函数, Crypt
// 加密因子M1, IA1, IC1看情况选取

/* test
data := make([]byte, 3)
data[0] = 97
data[1] = 98
data[2] = 99
_ := codec.Crypt(2, data)
fmt.Println(data)
=>result: [107 100 101]
*/

const (
	M1  uint32 = 1 << 19
	IA1 uint32 = 2 << 20
	IC1 uint32 = 3 << 21
)

func Crypt(key uint32, data []byte) (err error) {
	if key == 0 {
		key = 1
	}
	for i, _ := range data {
		key = IA1*(key%M1) + IC1
		data[i] ^= byte((key >> 20 & 0xff))
	}

	return nil
}
