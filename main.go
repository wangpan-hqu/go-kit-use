package main

import (
	"errors"
	"fmt"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	EndPoint1 "github.com/wangpan-hqu/go-kit-use/EndPoint"
	"github.com/wangpan-hqu/go-kit-use/Server"
	"github.com/wangpan-hqu/go-kit-use/Tool"
	"github.com/wangpan-hqu/go-kit-use/Transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1.先创建我们最开始定义的Server/server.go
	s := Server.Server{}

	// 2.在用EndPoint/endpoint.go 创建业务服务
	hello := EndPoint1.MakeServerEndPointHello(s)
	Bye := EndPoint1.MakeServerEndPointBye(s)

	// 3.使用 kit 创建 handler
	// 固定格式
	// 传入 业务服务 以及 定义的 加密解密方法
	helloServer := httpTransport.NewServer(hello, Transport.HelloDecodeRequest, Transport.HelloEncodeResponse)
	sayServer := httpTransport.NewServer(Bye, Transport.ByeDecodeRequest, Transport.ByeEncodeResponse)

	/*
		// 使用http包启动服务
		go http.ListenAndServe("0.0.0.0:8000", helloServer)

		go http.ListenAndServe("0.0.0.0:8001", sayServer)

		select {}
		*
	*/

	// https://github.com/gorilla/mux
	r := mux.NewRouter()
	// 注册路由
	r.Handle("/hello", helloServer)
	r.Handle("/bye", sayServer)

	// 因为这里要做服务发现,所以我们增加一个路由 进行心跳检测使用
	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	//	_ = http.ListenAndServe("0.0.0.0:8000", r)
	// 注册
	errChan := make(chan error)
	sign := make(chan os.Signal)
	go func() {
		err := Tool.RegService("127.0.0.1:8500", "1", "测试", "127.0.0.1", 8000, "5s", "http://127.0.0.1:8000/health", "test")
		if err != nil {
			errChan <- err
		}
		_ = http.ListenAndServe("0.0.0.0:8000", r)
	}()
	go func() {
		// 接收到信号
		signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
		<-sign
		errChan <- errors.New("0")
	}()
	fmt.Println(<-errChan)
	Tool.LogOutServer()
}
