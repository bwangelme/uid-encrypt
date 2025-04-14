# UID 加密库

这是一个高效的 UID 加密库，用于将 int64 类型的 UID 加密为定长字节序列，并支持解密还原。该库特别优化了对负数和小数字的处理，使用 ZigZag 编码来保持较小的输出大小。

## 特性

- 支持 int64 类型的 UID 加密和解密
- 使用 ZigZag 编码优化负数处理
- 对小数字（< UINT32_MAX）使用 16 字节输出
- 对大数字使用 24 字节输出
- 包含校验机制，可以检测数据完整性
- 使用 DES-CBC 加密，保证数据安全性
- 支持十六进制字符串编解码

## 安装

```bash
go get github.com/bwangelme/uid-encrypt
```

## 使用方法

```go
import "github.com/bwangelme/uid-encrypt/crypto"

// 加密 UID 为字节序列
uid := int64(12345)
ciphertext, err := crypto.EncryptUID(uid)
if err != nil {
    // 处理错误
}

// 解密字节序列为 UID
decryptedUID, err := crypto.DecryptUID(ciphertext)
if err != nil {
    // 处理错误
}

// 加密 UID 为十六进制字符串
hexStr, err := crypto.EncodeUID(uid)
if err != nil {
    // 处理错误
}

// 解密十六进制字符串为 UID
decryptedUID, err = crypto.DecodeUID(hexStr)
if err != nil {
    // 处理错误
}
```

## 输出大小

库会根据输入值的大小选择不同的输出格式：

1. 小数字格式（16 字节输出）：
   - 适用于 ZigZag 编码后小于 UINT32_MAX 的值
   - 包括大多数常见的 UID 值
   - 十六进制字符串长度：32 字符

2. 大数字格式（24 字节输出）：
   - 适用于 ZigZag 编码后大于等于 UINT32_MAX 的值
   - 十六进制字符串长度：48 字符

## ZigZag 编码

库使用 ZigZag 编码来优化负数的处理。编码规则如下：

- 0 -> 0
- -1 -> 1
- 1 -> 2
- -2 -> 3
- 2 -> 4
- ...

这种编码方式确保小的负数也能得到小的编码值，从而减少输出大小。

## 安全性

- 使用 DES-CBC 模式加密
- 包含魔数校验，可以检测数据完整性
- 包含冗余校验，可以检测数据损坏

## 限制

- DES 密钥和初始化向量是固定的
- 输出大小是固定的（16 字节或 24 字节）
- 不支持并发加密/解密操作的原子性

## 测试

```bash
go test -v ./...
```

## 性能测试结果

在 AMD Ryzen 7 3700X 8-Core Processor 上运行的性能测试结果：

### EncryptUID 函数

| 测试用例 | 执行时间 | 内存分配 | 分配次数 |
|---------|---------|---------|---------|
| 正数    | 3336 ns/op | 384 B/op | 16 allocs/op |
| 负数    | 3399 ns/op | 384 B/op | 16 allocs/op |
| 零      | 3251 ns/op | 384 B/op | 16 allocs/op |
| 大数    | 3248 ns/op | 384 B/op | 16 allocs/op |
| 大负数  | 3322 ns/op | 384 B/op | 16 allocs/op |

### EncodeUID 函数

| 测试用例 | 执行时间 | 内存分配 | 分配次数 |
|---------|---------|---------|---------|
| 正数    | 3564 ns/op | 448 B/op | 18 allocs/op |
| 负数    | 3625 ns/op | 448 B/op | 18 allocs/op |
| 零      | 3332 ns/op | 448 B/op | 18 allocs/op |
| 大数    | 3498 ns/op | 448 B/op | 18 allocs/op |
| 大负数  | 3619 ns/op | 448 B/op | 18 allocs/op |

### DecryptUID 函数

| 测试用例 | 执行时间 | 内存分配 | 分配次数 |
|---------|---------|---------|---------|
| 正数    | 3635 ns/op | 440 B/op | 18 allocs/op |
| 负数    | 3642 ns/op | 440 B/op | 19 allocs/op |
| 零      | 3358 ns/op | 424 B/op | 15 allocs/op |
| 大数    | 3567 ns/op | 440 B/op | 19 allocs/op |
| 大负数  | 3505 ns/op | 440 B/op | 19 allocs/op |

### DecodeUID 函数

| 测试用例 | 执行时间 | 内存分配 | 分配次数 |
|---------|---------|---------|---------|
| 正数    | 3622 ns/op | 456 B/op | 19 allocs/op |
| 负数    | 3741 ns/op | 456 B/op | 20 allocs/op |
| 零      | 3695 ns/op | 440 B/op | 16 allocs/op |
| 大数    | 3769 ns/op | 456 B/op | 20 allocs/op |
| 大负数  | 3626 ns/op | 456 B/op | 20 allocs/op |

性能测试命令：

```bash
make bench
```

## 注意事项

- 加密结果是一个十六进制字符串
- 支持的范围是 int64 的所有值
- 加密和解密操作都是对称的，同一个 UID 每次加密的结果都相同

## 许可证

MIT License
