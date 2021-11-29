package main

import (
	"context"
	flag "flag"
	"io"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
)

const (
	ListenAddr  = ":80"
	NormalLevel = 2
	DebugLevel  = 4
)

// 1.接收客户端 request，并将 request 中带的 header 写入 response header
// 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
// 4.当访问 localhost/healthz 时，应返回200
func main() {

	// 接收外部环境变量控制日志级别
	if os.Getenv("Debug") == "true" {
		flag.Set("v", "4")
	} else {
		flag.Set("v", "2")
	}

	// 默认就让日志输出在std，应用本身不负责日志采集
	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.V(NormalLevel).Info("Starting http server...")

	// BuildServerMux
	mux := buildServerMux()

	// Get ListenAddr
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ListenAddr
	}

	// Web Instance
	service := http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("listen: %s\n", err)
		}
	}()
	glog.V(NormalLevel).Info("Server Started")

	<-done

	glog.V(NormalLevel).Info("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := service.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Failed:%+v", err)
	}
	glog.V(NormalLevel).Info("Server Exited Properly")

}

// Build Server Mux
func buildServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	injectDebug(mux)
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthzHandler)

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
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok")
	t := time.Now().Format("2006-01-02 15:04:05")
	glog.V(DebugLevel).Infof("time: %v, health check \n", t)
}

// rootHandler Path /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// copy header form request
	for k, v := range r.Header {
		w.Header().Add(k, strings.Join(v, ","))
		glog.V(DebugLevel).Info(k, v)
	}
	vEnv := os.Getenv("VERSION")
	if vEnv == "" {
		vEnv = "0.0.0"
	}
	glog.V(DebugLevel).Info("vEnv: ", vEnv)
	w.Header().Add("VERSION", vEnv)
	w.WriteHeader(http.StatusOK)
	logResponse(http.StatusOK, r)
}

// print remote client address , port , status code
func logResponse(statusCode int, r *http.Request) {
	t := time.Now().Format("2006-01-02 15:04:05")
	glog.V(DebugLevel).Infof("time: %s, remote client: %s, status code: %d\n", t, r.RemoteAddr, statusCode)
}
