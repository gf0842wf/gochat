package zcodec

// 编解码常用函数

// 编码
func PutUInt8(num uint8, buf []byte, edian byte) {
	// len(buf) == 1
	buf[0] = byte(num)
}

func PutUInt16(num uint16, buf []byte, edian byte) {
	// len(buf) == 2
	buf[0] = byte(num >> 8)
	buf[1] = byte(num)
	if edian == 62 { // ">"

	} else if edian == 60 { // "<"
		buf[0] ^= buf[1]
		buf[1] ^= buf[0]
		buf[0] ^= buf[1]
	}
}

func PutUInt32(num uint32, buf []byte, edian byte) {
	// len(buf) == 4
	buf[0] = byte(num >> 24)
	buf[1] = byte(num >> 16)
	buf[2] = byte(num >> 8)
	buf[3] = byte(num)
	if edian == 62 {

	} else if edian == 60 {
		buf[0] ^= buf[3]
		buf[3] ^= buf[0]
		buf[0] ^= buf[3]

		buf[1] ^= buf[2]
		buf[2] ^= buf[1]
		buf[1] ^= buf[2]
	}
}

func PutUInt64(num uint64, buf []byte, edian byte) {
	// len(buf) == 8
	if edian == 62 {
		buf[0] = byte(num >> 56)
		buf[1] = byte(num >> 48)
		buf[2] = byte(num >> 40)
		buf[3] = byte(num >> 32)
		buf[4] = byte(num >> 24)
		buf[5] = byte(num >> 16)
		buf[6] = byte(num >> 8)
		buf[7] = byte(num)
	} else if edian == 60 {
		buf[0] = byte(num)
		buf[1] = byte(num >> 8)
		buf[2] = byte(num >> 16)
		buf[3] = byte(num >> 24)
		buf[4] = byte(num >> 32)
		buf[5] = byte(num >> 40)
		buf[6] = byte(num >> 48)
		buf[7] = byte(num >> 56)
	}
}

// 解码
func ToUInt8(buf []byte, edian byte) uint8 {
	// len(buf) == 1    -->B
	t := uint8(buf[0])
	return t
}

func ToUInt16(buf []byte, edian byte) uint16 {
	// len(buf) == 2    -->H
	t := uint16(buf[0])
	if edian == 62 { // ">"
		t = t<<8 | uint16(buf[1])
	} else if edian == 60 { // "<"
		t = t | uint16(buf[1])<<8
	}

	return t
}

func ToUInt32(buf []byte, edian byte) uint32 {
	// len(buf) == 4    -->I
	t := uint32(buf[0])
	if edian == 62 {
		t = t << 24
		t = t | uint32(buf[1])<<16
		t = t | uint32(buf[2])<<8
		t = t | uint32(buf[3])

	} else if edian == 60 {
		t = t | uint32(buf[1])<<8
		t = t | uint32(buf[2])<<16
		t = t | uint32(buf[3])<<24
	}
	return t
}

func ToUInt64(buf []byte, edian byte) uint64 {
	//len(buf) == 8    -->Q
	t := uint64(buf[0])
	if edian == 62 {
		t = t << 56
		t = t | uint64(buf[1])<<48
		t = t | uint64(buf[2])<<40
		t = t | uint64(buf[3])<<32
		t = t | uint64(buf[4])<<24
		t = t | uint64(buf[5])<<16
		t = t | uint64(buf[6])<<8
		t = t | uint64(buf[7])
	} else if edian == 60 {
		t = t | uint64(buf[1])<<8
		t = t | uint64(buf[2])<<16
		t = t | uint64(buf[3])<<24
		t = t | uint64(buf[4])<<32
		t = t | uint64(buf[5])<<40
		t = t | uint64(buf[6])<<48
		t = t | uint64(buf[7])<<56
	}
	return t
}
