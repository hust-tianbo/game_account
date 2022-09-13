package main

import (
	"fmt"
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
		w.Write([]byte("welcome to account"))
	})
	return mux
}
