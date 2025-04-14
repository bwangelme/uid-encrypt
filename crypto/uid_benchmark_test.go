package crypto

import (
	"testing"
)

func BenchmarkEncryptUID(b *testing.B) {
	testCases := []struct {
		name string
		uid  int64
	}{
		{"正数", 12345},
		{"负数", -12345},
		{"零", 0},
		{"大数", 999999999},
		{"大负数", -999999999},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := EncryptUID(tc.uid)
				if err != nil {
					b.Fatalf("加密失败: %v", err)
				}
			}
		})
	}
}

func BenchmarkEncodeUID(b *testing.B) {
	testCases := []struct {
		name string
		uid  int64
	}{
		{"正数", 12345},
		{"负数", -12345},
		{"零", 0},
		{"大数", 999999999},
		{"大负数", -999999999},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := EncodeUID(tc.uid)
				if err != nil {
					b.Fatalf("编码失败: %v", err)
				}
			}
		})
	}
}

func BenchmarkDecryptUID(b *testing.B) {
	// 准备测试数据
	testCases := []struct {
		name       string
		uid        int64
		ciphertext []byte
	}{
		{"正数", 12345, nil},
		{"负数", -12345, nil},
		{"零", 0, nil},
		{"大数", 999999999, nil},
		{"大负数", -999999999, nil},
	}

	// 预先加密所有测试数据
	for i := range testCases {
		ciphertext, err := EncryptUID(testCases[i].uid)
		if err != nil {
			b.Fatalf("准备测试数据失败: %v", err)
		}
		testCases[i].ciphertext = ciphertext
	}

	// 运行性能测试
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := DecryptUID(tc.ciphertext)
				if err != nil {
					b.Fatalf("解密失败: %v", err)
				}
			}
		})
	}
}

func BenchmarkDecodeUID(b *testing.B) {
	// 准备测试数据
	testCases := []struct {
		name   string
		uid    int64
		hexStr string
	}{
		{"正数", 12345, ""},
		{"负数", -12345, ""},
		{"零", 0, ""},
		{"大数", 999999999, ""},
		{"大负数", -999999999, ""},
	}

	// 预先编码所有测试数据
	for i := range testCases {
		hexStr, err := EncodeUID(testCases[i].uid)
		if err != nil {
			b.Fatalf("准备测试数据失败: %v", err)
		}
		testCases[i].hexStr = hexStr
	}

	// 运行性能测试
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := DecodeUID(tc.hexStr)
				if err != nil {
					b.Fatalf("解码失败: %v", err)
				}
			}
		})
	}
}
