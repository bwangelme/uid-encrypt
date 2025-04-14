package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/bwangelme/uid-encrypt/crypto"
)

func main() {
	// 定义子命令
	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)

	// 检查参数
	if len(os.Args) < 2 {
		fmt.Println("请指定子命令: encrypt 或 decrypt")
		fmt.Println("\n用法:")
		fmt.Println("  uid encrypt <uid>     - 加密一个 UID")
		fmt.Println("  uid decrypt <hexstr>  - 解密一个十六进制字符串")
		fmt.Println("\n示例:")
		fmt.Println("  uid encrypt 12345")
		fmt.Println("  uid encrypt -1")
		fmt.Println("  uid decrypt b2dac66902d57757e4d1d531fad89ad5")
		os.Exit(1)
	}

	// 根据子命令执行不同的操作
	switch os.Args[1] {
	case "encrypt":
		encryptCmd.Parse(os.Args[2:])
		if encryptCmd.NArg() != 1 {
			fmt.Println("请提供要加密的 UID")
			os.Exit(1)
		}
		// 解析 UID
		uid, err := strconv.ParseInt(encryptCmd.Arg(0), 10, 64)
		if err != nil {
			fmt.Printf("无效的 UID: %v\n", err)
			os.Exit(1)
		}
		// 加密
		hexStr, err := crypto.EncodeUID(uid)
		if err != nil {
			fmt.Printf("加密失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("加密结果: %s\n", hexStr)

	case "decrypt":
		decryptCmd.Parse(os.Args[2:])
		if decryptCmd.NArg() != 1 {
			fmt.Println("请提供要解密的十六进制字符串")
			os.Exit(1)
		}
		// 解密
		uid, err := crypto.DecodeUID(decryptCmd.Arg(0))
		if err != nil {
			fmt.Printf("解密失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("解密结果: %d\n", uid)

	default:
		fmt.Printf("未知的子命令: %s\n", os.Args[1])
		os.Exit(1)
	}
}
