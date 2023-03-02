package main

import (
	"Morse-shellcode/server/tools"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/socket", tools.SocketHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
