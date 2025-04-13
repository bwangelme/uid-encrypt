package crypto

import (
	"testing"
)

func TestZigZagEncodeDecode(t *testing.T) {
	// 测试 -1 的编码
	value := int64(-1)
	encoded := zigzagEncode(value)
	decoded := zigzagDecode(encoded)

	if encoded != 1 {
		t.Errorf("zigzagEncode(-1) = %d; 期望 1", encoded)
	}
	if decoded != -1 {
		t.Errorf("zigzagDecode(1) = %d; 期望 -1", decoded)
	}

	// 测试更多值
	testValues := []int64{-10, -5, -2, -1, 0, 1, 2, 5, 10}
	expectedEncoded := []uint64{19, 9, 3, 1, 0, 2, 4, 10, 20}

	for i, v := range testValues {
		enc := zigzagEncode(v)
		dec := zigzagDecode(enc)

		if enc != expectedEncoded[i] {
			t.Errorf("zigzagEncode(%d) = %d; 期望 %d", v, enc, expectedEncoded[i])
		}
		if dec != v {
			t.Errorf("zigzagDecode(%d) = %d; 期望 %d", enc, dec, v)
		}
	}
}
