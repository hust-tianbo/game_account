package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hust-tianbo/game_account/internal/logic"
	"github.com/hust-tianbo/go_lib/log"
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
		body, _ := ioutil.ReadAll(r.Body)
		var req logic.CheckAuthReq
		json.Unmarshal(body, &req)
		var rsp logic.CheckAuthRsp
		defer func() {
			log.Debugf("[GetHttpServerMux]deal log:%+v,%+v", req, rsp)
		}()

		rsp = logic.CheckAuth(req)
		resBytes, _ := json.Marshal(rsp)
		w.Write([]byte(resBytes))
	})
	return mux
}
