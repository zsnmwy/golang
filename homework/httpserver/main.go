package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	ListenAddr  = ":80"
	NormalLevel = 2
)

// 1.接收客户端 request，并将 request 中带的 header 写入 response header
// 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
// 4.当访问 localhost/healthz 时，应返回200
func main() {

	glog.V(NormalLevel).Info("Starting http server...")

	// BuildServerMux
	mux := buildServerMux()

	// Get ListenAddr
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ListenAddr
	}

	// Boot web application
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}

}

// Build Server Mux
func buildServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	injectDebug(mux)
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)

	return mux
}

// make application easy to debug
func injectDebug(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// check application is health
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok")
}

// rootHandler Path /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// copy header form request
	for k, v := range r.Header {
		w.Header().Add(k, strings.Join(v, ","))
		fmt.Println(k, v)
	}
	vEnv := os.Getenv("VERSION")
	if vEnv == "" {
		vEnv = "0.0.0"
	}
	w.Header().Add("VERSION", vEnv)
	w.WriteHeader(http.StatusOK)
	logResponse(http.StatusOK, r)
}

// print remote client address , port , status code
func logResponse(statusCode int, r *http.Request) {
	t := time.Now().Format("2006-01-02 15-04-05")
	fmt.Printf("time: %s, remote client: %s, status code: %d", t, r.RemoteAddr, statusCode)
}
