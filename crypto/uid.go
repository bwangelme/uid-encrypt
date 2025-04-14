package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
)

const (
	magicNum1 = 12345      // 第一个魔数，需要小于 uint16max
	magicNum2 = 54321      // 第二个魔数，需要小于 uint16max
	desKey    = "12345678" // DES 密钥，8字节
	desIV     = "87654321" // DES 初始化向量，8字节
)

// pack 根据format字符串打包数字
// format 格式说明：
// N - 大端写入 uint32
// n - 大端写入 uint16
// C - 写入 uint8
// V - 小端写入 uint32
// v - 小端写入 uint16
func pack(format string, numbers ...interface{}) ([]byte, error) {
	if len(format) != len(numbers) {
		return nil, errors.New("format长度与参数数量不匹配")
	}

	buf := new(bytes.Buffer)
	for i, f := range format {
		switch f {
		case 'N':
			if num, ok := numbers[i].(uint32); ok {
				if err := binary.Write(buf, binary.BigEndian, num); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("参数类型不匹配: 需要uint32")
			}
		case 'n':
			if num, ok := numbers[i].(uint16); ok {
				if err := binary.Write(buf, binary.BigEndian, num); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("参数类型不匹配: 需要uint16")
			}
		case 'C':
			if num, ok := numbers[i].(uint8); ok {
				if err := binary.Write(buf, binary.BigEndian, num); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("参数类型不匹配: 需要uint8")
			}
		case 'V':
			if num, ok := numbers[i].(uint32); ok {
				if err := binary.Write(buf, binary.LittleEndian, num); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("参数类型不匹配: 需要uint32")
			}
		case 'v':
			if num, ok := numbers[i].(uint16); ok {
				if err := binary.Write(buf, binary.LittleEndian, num); err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("参数类型不匹配: 需要uint16")
			}
		default:
			return nil, errors.New("不支持的格式字符: " + string(f))
		}
	}

	return buf.Bytes(), nil
}

// unpack 根据format字符串解包数据
// format 格式说明：
// N - 大端读取 uint32
// n - 大端读取 uint16
// C - 读取 uint8
// V - 小端读取 uint32
// v - 小端读取 uint16
func unpack(format string, data []byte) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("数据为空")
	}

	buf := bytes.NewReader(data)
	result := make([]interface{}, 0, len(format))

	for _, f := range format {
		switch f {
		case 'N':
			var num uint32
			if err := binary.Read(buf, binary.BigEndian, &num); err != nil {
				return nil, fmt.Errorf("读取大端uint32失败: %v", err)
			}
			result = append(result, num)
		case 'n':
			var num uint16
			if err := binary.Read(buf, binary.BigEndian, &num); err != nil {
				return nil, fmt.Errorf("读取大端uint16失败: %v", err)
			}
			result = append(result, num)
		case 'C':
			var num uint8
			if err := binary.Read(buf, binary.BigEndian, &num); err != nil {
				return nil, fmt.Errorf("读取uint8失败: %v", err)
			}
			result = append(result, num)
		case 'V':
			var num uint32
			if err := binary.Read(buf, binary.LittleEndian, &num); err != nil {
				return nil, fmt.Errorf("读取小端uint32失败: %v", err)
			}
			result = append(result, num)
		case 'v':
			var num uint16
			if err := binary.Read(buf, binary.LittleEndian, &num); err != nil {
				return nil, fmt.Errorf("读取小端uint16失败: %v", err)
			}
			result = append(result, num)
		default:
			return nil, fmt.Errorf("不支持的格式字符: %c", f)
		}
	}

	// 检查是否还有剩余数据
	if buf.Len() > 0 {
		return nil, errors.New("数据长度超过格式要求")
	}

	return result, nil
}

// zigzagEncode 将 int64 转换为无符号整数
func zigzagEncode(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

// zigzagDecode 将无符号整数转换回 int64
func zigzagDecode(n uint64) int64 {
	return int64(n>>1) ^ -int64(n&1)
}

// EncryptUID 将 int64 类型的 UID 加密为 16 字节的密文
func EncryptUID(uid int64) ([]byte, error) {
	var data []byte
	var err error

	// 使用 ZigZag 编码将负数转换为无符号整数
	zigzagUid := zigzagEncode(uid)

	// 检查编码后的 uid 是否小于 uint32 最大值
	if zigzagUid < math.MaxUint32 {
		// 使用 pack 函数打包数据：N(大端uint32) + n(大端uint16) + n(大端uint16) + V(小端int64) + v(小端uint16) + C(uint8)
		data, err = pack("NnCCVvC",
			uint32(zigzagUid),           // 大端写入 zigzag 编码后的 uid
			uint16(zigzagUid%magicNum1), // 大端写入 zigzag 编码后的 uid % magicNum1
			uint8(0),
			uint8(0),
			uint32(zigzagUid),           // 小端写入 zigzag 编码后的 uid
			uint16(zigzagUid%magicNum2), // 小端写入 zigzag 编码后的 uid % magicNum2
			uint8(0),                    // 0
		)
	} else {
		// 对于大于 UINT32_MAX 的数字，使用 24 字节序列化
		// 格式：NNnCCVVvC
		// N: zigzagUid>>32 (高32位)
		// N: zigzagUid%UINT32_MAX (低32位)
		// n: zigzagUid%magicNum1
		// C: 0
		// C: 0
		// V: zigzagUid>>32 (高32位)
		// V: zigzagUid%UINT32_MAX (低32位)
		// v: zigzagUid%magicNum2
		// C: 0
		data, err = pack("NNnCCVVvC",
			uint32(zigzagUid>>32),       // 高32位
			uint32(zigzagUid),           // 低32位
			uint16(zigzagUid%magicNum1), // zigzag 编码后的 uid % magicNum1
			uint8(0),                    // 0
			uint8(0),                    // 0
			uint32(zigzagUid>>32),       // 高32位
			uint32(zigzagUid),           // 低32位
			uint16(zigzagUid%magicNum2), // zigzag 编码后的 uid % magicNum2
			uint8(0),                    // 0
		)
	}

	if err != nil {
		return nil, err
	}

	// 确保数据长度为 15 或 23 字节
	if len(data) != 15 && len(data) != 23 {
		return nil, errors.New("数据长度不正确")
	}

	// 进行 DES 加密
	block, err := des.NewCipher([]byte(desKey))
	if err != nil {
		return nil, err
	}

	// 使用 CBC 模式加密
	iv := []byte(desIV)
	mode := cipher.NewCBCEncrypter(block, iv)

	// 确保数据长度是 8 的倍数（DES 块大小）
	paddedData := pkcs7Padding(data)

	// 加密数据
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	return ciphertext, nil
}

// EncodeUID 将 int64 类型的 UID 加密并转换为十六进制字符串
func EncodeUID(uid int64) (string, error) {
	// 先加密
	ciphertext, err := EncryptUID(uid)
	if err != nil {
		return "", err
	}

	// 确保密文长度是块大小的整数倍
	if len(ciphertext)%8 != 0 {
		return "", errors.New("加密后的密文长度不是块大小的整数倍")
	}

	// 转换为十六进制字符串，每个字节输出两个字符
	return hex.EncodeToString(ciphertext), nil
}

// DecryptUID 将加密后的字节数组解密回 int64 类型的 UID
func DecryptUID(ciphertext []byte) (int64, error) {
	// 检查输入长度
	if len(ciphertext) == 0 {
		return 0, errors.New("密文为空")
	}

	// 确保密文长度是块大小的整数倍
	if len(ciphertext)%des.BlockSize != 0 {
		return 0, errors.New("密文长度不是块大小的整数倍")
	}

	// 创建 DES 解密器
	block, err := des.NewCipher([]byte(desKey))
	if err != nil {
		return 0, err
	}

	// 使用 CBC 模式解密
	iv := []byte(desIV)
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	unpaddedData, err := pkcs7Unpadding(plaintext)
	if err != nil {
		return 0, fmt.Errorf("去除填充失败: %v", err)
	}

	// 根据长度判断是大数字还是小数字
	if len(unpaddedData) == 15 {
		// 15字节格式：N(uid) + n(uid%magicNum1) + n(0) + V(uid) + v(uid%magicNum2) + C(0)
		values, err := unpack("NnCCVvC", unpaddedData)
		if err != nil {
			return 0, fmt.Errorf("解包数据失败: %v", err)
		}

		// 获取解包后的值
		zigzagUid := values[0].(uint32)
		checkNum1 := values[1].(uint16)
		_ = values[2].(uint8) // 跳过0
		_ = values[3].(uint8) // 跳过0
		zigzagUid2 := values[4].(uint32)
		checkNum2 := values[5].(uint16)
		_ = values[6].(uint8) // 跳过0

		// 验证校验值
		if uint16(zigzagUid%magicNum1) != checkNum1 {
			return 0, errors.New("校验值1不匹配")
		}
		if uint16(zigzagUid2%magicNum2) != checkNum2 {
			return 0, errors.New("校验值2不匹配")
		}

		// 验证两个uid是否相同
		if zigzagUid != zigzagUid2 {
			return 0, errors.New("uid不匹配")
		}

		// 使用 ZigZag 解码将无符号整数转换回 int64
		return zigzagDecode(uint64(zigzagUid)), nil
	} else if len(unpaddedData) == 23 {
		// 23字节格式：N(zigzagUid>>32) + N(zigzagUid%UINT32_MAX) + n(zigzagUid%magicNum1) + C(0) + C(0) + V(zigzagUid>>32) + V(zigzagUid%UINT32_MAX) + v(zigzagUid%magicNum2) + C(0)
		values, err := unpack("NNnCCVVvC", unpaddedData)
		if err != nil {
			return 0, fmt.Errorf("解包数据失败: %v", err)
		}

		// 获取解包后的值
		highBits := values[0].(uint32)
		lowBits := values[1].(uint32)
		checkNum1 := values[2].(uint16)
		_ = values[3].(uint8) // 跳过0
		_ = values[4].(uint8) // 跳过0
		highBits2 := values[5].(uint32)
		lowBits2 := values[6].(uint32)
		checkNum2 := values[7].(uint16)
		_ = values[8].(uint8) // 跳过0

		// 重构 ZigZag 编码后的 uid
		zigzagUid := (uint64(highBits) << 32) | uint64(lowBits)

		// 验证校验值
		if uint16(zigzagUid%magicNum1) != checkNum1 {
			return 0, errors.New("校验值1不匹配")
		}
		if uint16(zigzagUid%magicNum2) != checkNum2 {
			return 0, errors.New("校验值2不匹配")
		}

		// 验证高32位和低32位是否匹配
		if highBits != highBits2 || lowBits != lowBits2 {
			return 0, errors.New("uid部分不匹配")
		}

		// 使用 ZigZag 解码将无符号整数转换回 int64
		return zigzagDecode(zigzagUid), nil
	} else {
		return 0, fmt.Errorf("解密后的数据长度不正确: %d", len(unpaddedData))
	}
}

// DecodeUID 将十六进制字符串解密回 int64 类型的 UID
func DecodeUID(hexStr string) (int64, error) {
	// 检查输入长度
	if len(hexStr) == 0 {
		return 0, errors.New("十六进制字符串为空")
	}

	// 将十六进制字符串转换为字节数组
	ciphertext, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, fmt.Errorf("无效的十六进制字符串: %v", err)
	}

	// 检查密文长度
	if len(ciphertext) == 0 {
		return 0, errors.New("解码后的密文为空")
	}

	// 确保密文长度是块大小的整数倍
	if len(ciphertext)%8 != 0 {
		return 0, errors.New("解码后的密文长度不是块大小的整数倍")
	}

	return DecryptUID(ciphertext)
}

// pkcs7Padding 对数据进行 PKCS7 填充
func pkcs7Padding(data []byte) []byte {
	blockSize := 8 // DES 块大小
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpadding 去除 PKCS7 填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("数据为空")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("填充长度不正确")
	}

	// 验证所有填充字节是否相同
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("填充字节不一致")
		}
	}

	return data[:length-padding], nil
}
