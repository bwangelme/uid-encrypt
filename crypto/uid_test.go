package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeUID(t *testing.T) {
	tests := []struct {
		name     string
		uid      int64
		wantLen  int
		expected string // 期望的十六进制字符串前缀（因为加密结果每次可能不同，我们只检查前缀）
	}{
		{
			name:     "小数字测试-12345",
			uid:      12345,
			wantLen:  32,
			expected: "b2dac66902d57757e4d1d531fad89ad5",
		},
		{
			name:     "大数字测试-UINT32_MAX+1",
			uid:      4294967296,
			wantLen:  48,
			expected: "f6cbb4d381d596db43c560f7c19c81498d21ba0089d19117",
		},
		{
			name:     "零值测试",
			uid:      0,
			wantLen:  32,
			expected: "385430289b759424846e483e21fb257e",
		},
		{
			name:     "负数测试",
			uid:      -12345,
			wantLen:  32,
			expected: "86c9615314930abeaab421e2cfc08ec4",
		},
		{
			name:     "边界值测试-最大16字节数",
			uid:      2147483647,
			wantLen:  32,
			expected: "045b89374f05b126b6b42b6150f1f64a",
		},
		{
			name:     "边界值测试-UINT32_MAX",
			uid:      4294967295,
			wantLen:  48,
			expected: "9de416d9aeb35d0c73df74d22b8c9bfea6ff89a8eb780a85",
		},
		{
			name:     "边界值测试-UINT32_MAX+1",
			uid:      4294967296,
			wantLen:  48,
			expected: "f6cbb4d381d596db43c560f7c19c81498d21ba0089d19117",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeUID(tt.uid)
			assert.NoError(t, err, "EncodeUID() 不应该返回错误")
			assert.NotEmpty(t, got, "EncodeUID() 不应该返回空字符串")
			assert.Equal(t, tt.wantLen, len(got), "EncodeUID() 返回的字符串长度不正确")

			// 检查是否都是十六进制字符
			for i, c := range got {
				assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'),
					"位置 %d 的字符 %c 不是十六进制字符", i, c)
			}

			// 如果指定了期望值，检查前缀
			if tt.expected != "" {
				assert.True(t, len(got) >= len(tt.expected), "结果长度应该大于等于期望值长度")
				assert.Equal(t, tt.expected, got[:len(tt.expected)], "结果前缀不匹配")
			}
		})
	}
}

func TestEncryptUID(t *testing.T) {
	tests := []struct {
		name     string
		uid      int64
		wantLen  int
		expected []byte // 期望的字节序列（可选）
	}{
		{
			name:    "小数字测试/12345",
			uid:     12345,
			wantLen: 16,
		},
		{
			name:    "大数字测试-UINT32_MAX+1",
			uid:     4294967296,
			wantLen: 24,
		},
		{
			name:    "边界值测试-UINT32_MAX",
			uid:     4294967295,
			wantLen: 24,
		},
		{
			name:    "边界值测试-UINT32_MAX+1",
			uid:     4294967296,
			wantLen: 24,
		},
		{
			name:    "特殊值测试-1",
			uid:     1,
			wantLen: 16,
		},
		{
			name:    "特殊值测试-(-1)",
			uid:     -1,
			wantLen: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptUID(tt.uid)
			assert.NoError(t, err, "EncryptUID() 不应该返回错误")
			assert.NotEmpty(t, got, "EncryptUID() 不应该返回空字节数组")
			assert.Equal(t, tt.wantLen, len(got), "EncryptUID() 返回的字节长度不正确")

			// 检查是否所有字节都不为0
			allZero := true
			for _, b := range got {
				if b != 0 {
					allZero = false
					break
				}
			}
			assert.False(t, allZero, "EncryptUID() 加密结果不应该全为0")

			// 如果指定了期望值，检查字节序列
			if tt.expected != nil {
				assert.Equal(t, tt.expected, got, "加密结果与期望值不匹配")
			}

			// 检查相同输入是否产生相同的输出（DES加密应该是确定性的）
			got2, err := EncryptUID(tt.uid)
			assert.NoError(t, err, "第二次加密不应该返回错误")
			assert.Equal(t, got, got2, "相同输入应该产生相同的输出")
		})
	}
}

func TestDecryptUID(t *testing.T) {
	tests := []struct {
		name      string
		uid       int64
		wantBytes int
	}{
		{
			name:      "小数字测试/12345",
			uid:       12345,
			wantBytes: 16,
		},
		{
			name:      "大数字测试/UINT32_MAX+3",
			uid:       4294967298,
			wantBytes: 24,
		},
		{
			name:      "边界值测试/UINT32_MAX",
			uid:       4294967295,
			wantBytes: 24,
		},
		{
			name:      "边界值测试/UINT32_MAX+1",
			uid:       4294967296,
			wantBytes: 24,
		},
		{
			name:      "特殊值测试/1",
			uid:       1,
			wantBytes: 16,
		},
		{
			name:      "特殊值测试/(-1)",
			uid:       -1,
			wantBytes: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 先加密
			ciphertext, err := EncryptUID(tt.uid)
			assert.NoError(t, err, "加密不应该返回错误")
			assert.Equal(t, tt.wantBytes, len(ciphertext), "加密结果长度不正确")

			// 再解密
			got, err := DecryptUID(ciphertext)
			assert.NoError(t, err, "解密不应该返回错误")
			assert.Equal(t, tt.uid, got, "解密结果与原始值不匹配")

			// 测试错误情况
			// 1. 错误的长度
			invalidCiphertext := make([]byte, tt.wantBytes+1)
			_, err = DecryptUID(invalidCiphertext)
			assert.Error(t, err, "错误的长度应该返回错误")

			// 2. 错误的校验值
			if tt.wantBytes == 16 {
				// 修改第一个校验值
				ciphertext[4] ^= 0xFF
				_, err = DecryptUID(ciphertext)
				assert.Error(t, err, "错误的校验值应该返回错误")
			}
		})
	}
}

func TestDecodeUID(t *testing.T) {
	tests := []struct {
		name    string
		uid     int64
		wantLen int
	}{
		{
			name:    "小数字测试/12345",
			uid:     12345,
			wantLen: 32, // 16字节 * 2（每个字节两个十六进制字符）
		},
		{
			name:    "大数字测试/UINT32_MAX+3",
			uid:     4294967298,
			wantLen: 48, // 24字节 * 2（每个字节两个十六进制字符）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 先编码
			hexStr, err := EncodeUID(tt.uid)
			assert.NoError(t, err, "编码不应该返回错误")
			assert.Equal(t, tt.wantLen, len(hexStr), "编码结果长度不正确")

			// 再解码
			got, err := DecodeUID(hexStr)
			assert.NoError(t, err, "解码不应该返回错误")
			assert.Equal(t, tt.uid, got, "解码结果与原始值不匹配")

			// 测试错误情况
			// 1. 错误的长度
			_, err = DecodeUID(hexStr[:len(hexStr)-1])
			assert.Error(t, err, "错误的长度应该返回错误")

			// 2. 无效的十六进制字符串
			_, err = DecodeUID("invalid hex string")
			assert.Error(t, err, "无效的十六进制字符串应该返回错误")
		})
	}
}
