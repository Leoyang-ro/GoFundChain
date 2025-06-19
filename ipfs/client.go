# IPFS 上传/下载

package main

import (
	"fmt"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	// 创建 IPFS Shell，连接到本地节点
	sh := shell.NewShell("localhost:5001")

	// 方式 1：上传文件内容（字符串）
	content := "Hello, IPFS from Golang!"
	cid, err := sh.Add(strings.NewReader(content))
	if err != nil {
		fmt.Println("上传失败:", err)
		return
	}
	fmt.Println("字符串内容上传成功，CID:", cid)

	// 方式 2：上传本地文件
	file, err := os.Open("testfile.txt")
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close()

	fileCID, err := sh.Add(file)
	if err != nil {
		fmt.Println("文件上传失败:", err)
		return
	}
	fmt.Println("文件上传成功，CID:", fileCID)

	// 可访问：https://ipfs.io/ipfs/<CID>
}