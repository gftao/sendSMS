package comm

import (
	"fmt"
	"mygolib/modules/config"
	"mygolib/modules/myLogger"
	"errors"
	"time"
	"os"
	"os/signal"
	"net/http"
	"log"
	"context"
	"syscall"
)

type httpSvrConf struct {
	HttpsCertFile string
	HttpsKeyFile  string
	ListenIp      string
	ListenPort    int
	RecvTimeOut   int
	WriteTimeOut  int
	MaxAccNum     int
}

type HttpSvr struct {
	conf *httpSvrConf
}

func (t *HttpSvr) InitConfig() error {

	if !config.HasConfigInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}
	if !myLogger.HasLoggerInit() {
		return errors.New("日志模块未初始化，请先初始化")
	}

	config.SetSection("server")

	cf := &httpSvrConf{}
	cf.HttpsCertFile = config.StringDefault("HttpsCertFile", "")
	cf.HttpsKeyFile = config.StringDefault("HttpsKeyFile", "")
	cf.ListenIp = config.StringDefault("host", "")
	cf.ListenPort = config.IntDefault("port", 9090)
	cf.RecvTimeOut = config.IntDefault("readTimeout", 30)
	cf.WriteTimeOut = config.IntDefault("writeTimeout", 30)
	t.conf = cf
	if cf.ListenIp == "" || cf.ListenPort == 0 {
		return errors.New("https初始化失败，参数非法")
	}
	fmt.Println("HttpSvr加载成功")

	return nil
}
func (t *HttpSvr) RunSvr(h http.Handler) {
	srv := &http.Server{
		Addr:           t.conf.ListenIp + fmt.Sprintf(":%d", t.conf.ListenPort),
		Handler:        h,
		ReadTimeout:    time.Duration(t.conf.RecvTimeOut) * time.Second,
		WriteTimeout:   time.Duration(t.conf.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		defer myLogger.Info("----HttpSvr closed----")

		if t.conf.HttpsCertFile != "" && t.conf.HttpsKeyFile != "" {
			myLogger.Info("----HttpsSvr started----")
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		} else {
			myLogger.Info("----HttpSvr started----")
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
}
