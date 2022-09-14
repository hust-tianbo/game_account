package main

import (
	"encoding/json"
	"fmt"
	"github.com/hust-tianbo/game_account/internal"
	"net/http"
)

const (
	port = ":50052"
)

func main() {
	//modbus.get()
	fmt.Printf("begin account server")

	// 注册http接口
	mux := GetHttpServerMux()
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}),
	)

	fmt.Printf("end account server")
}

func GetHttpServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := internal.CheckAuth()
		resBytes, _ := json.Marshal(res)
		w.Write([]byte(resBytes))
	})
	return mux
}
