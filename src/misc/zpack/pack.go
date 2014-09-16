package zpack

// 类似python struct的编解码
// s: ">HIB12BQ"
// s: "<HIB12BQ"

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func IsPackNum(n byte) bool {
	return (n >= byte('0') && n <= byte('9'))
}

func CalcSize(s string) (size int) {
	// 格式只需支持 BHIQ 这4种
	// 状态机: 0->0 | 0->1->2->0

	s = s[1:]

	times := 1
	pos1, pos2 := 0, 0
	status := 0

	for i, v := range s {
		switch v {
		case 'B':
			if status == 0 {
				size += 1
			} else if status == 2 {
				size += (times * 1)
				status = 0
			}
		case 'H':
			if status == 0 {
				size += 2
			} else if status == 2 {
				size += (times * 2)
				status = 0
			}
		case 'I':
			if status == 0 {
				size += 4
			} else if status == 2 {
				size += (times * 4)
				status = 0
			}
		case 'Q':
			if status == 0 {
				size += 8
			} else if status == 2 {
				size += (times * 8)
				status = 0
			}
		default: //  只支持上面5种,剩下的是数字
			if status == 0 {
				pos1 = i
				status = 1
			}
			if !IsPackNum(byte(s[i+1])) {
				pos2 = i
				times, _ = strconv.Atoi(s[pos1 : pos2+1])
				status = 2
			}
		}

	}

	return
}

func Unpack(s string, data []byte) (msg []interface{}) {
	// ">2B" => msg: [[3,3]],  ">BB" => msg:[3,3]
	dian := s[0]
	count := 0 // 不是字节数,是元素个数
	s = s[1:]
	buf := bytes.NewBuffer(data)

	for _, v := range s {
		if !IsPackNum(byte(v)) {
			count += 1
		}
	}

	msg = make([]interface{}, count)

	times := 1
	pos1, pos2 := 0, 0
	status := 0

	j := 0

	for i, v := range s {
		switch v {
		case 'B':
			if status == 0 {
				msg[j], _ = buf.ReadByte()
			} else if status == 2 {
				m := make([]byte, times)
				for k, _ := range m {
					m[k], _ = buf.ReadByte()
				}
				msg[j] = m
				status = 0
			}
			j++
		case 'H':
			if status == 0 {
				var m uint16
				if dian == '>' {
					binary.Read(buf, binary.BigEndian, &m)
				} else if dian == '<' {
					binary.Read(buf, binary.LittleEndian, &m)
				}
				msg[j] = m
			} else if status == 2 {
				ms := make([]uint16, times)
				for k, _ := range ms {
					if dian == '>' {
						binary.Read(buf, binary.BigEndian, &ms[k])
					} else if dian == '<' {
						binary.Read(buf, binary.LittleEndian, &ms[k])
					}
				}
				msg[j] = ms
				status = 0
			}
			j++
		case 'I':
			if status == 0 {
				var m uint32
				if dian == '>' {
					binary.Read(buf, binary.BigEndian, &m)
				} else if dian == '<' {
					binary.Read(buf, binary.LittleEndian, &m)
				}
				msg[j] = m
			} else if status == 2 {
				ms := make([]uint32, times)
				for k, _ := range ms {
					if dian == '>' {
						binary.Read(buf, binary.BigEndian, &ms[k])
					} else if dian == '<' {
						binary.Read(buf, binary.LittleEndian, &ms[k])
					}
				}
				msg[j] = ms
				status = 0
			}
			j++
		case 'Q':
			if status == 0 {
				var m uint64
				if dian == '>' {
					binary.Read(buf, binary.BigEndian, &m)
				} else if dian == '<' {
					binary.Read(buf, binary.LittleEndian, &m)
				}
				msg[j] = m
			} else if status == 2 {
				ms := make([]uint64, times)
				for k, _ := range ms {
					if dian == '>' {
						binary.Read(buf, binary.BigEndian, &ms[k])
					} else if dian == '<' {
						binary.Read(buf, binary.LittleEndian, &ms[k])
					}
				}
				msg[j] = ms
				status = 0
			}
			j++
		default: //  只支持上面5种,剩下的是数字
			if status == 0 {
				pos1 = i
				status = 1
			}
			if !IsPackNum(byte(s[i+1])) {
				pos2 = i
				times, _ = strconv.Atoi(s[pos1 : pos2+1])
				status = 2
			}
		}
	}

	return
}

func Pack(dian byte, msg []interface{}) (data []byte) {
	buf := new(bytes.Buffer)
	for _, v := range msg {
		if dian == byte('>') {
			binary.Write(buf, binary.BigEndian, v)
		} else if dian == byte('<') {
			binary.Write(buf, binary.LittleEndian, v)
		}
	}
	data = buf.Bytes()
	return
}

//func main() {
//	s := ">B2HIQ"
//	fmt.Println(CalcSize(s))
//	data := []byte{4, 4, 188, 4, 188, 0, 0, 130, 53, 0, 0, 8, 21, 155, 16, 142, 56}
//	msg := Unpack(s, data)
//	fmt.Println(msg)
//	fmt.Println(Pack('>', msg))
//	//17
//	//[4 [1212 1212] 33333 8888888888888]
//	//[4 4 188 4 188 0 0 130 53 0 0 8 21 155 16 142 56]
//}
