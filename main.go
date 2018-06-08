package main

import (
	"runtime"
	"flag"
	"fmt"
	"os"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"mygolib/modules/myLogger"
	"mygolib/modules/config"
	"sendSMS/handle"
	"sendSMS/trans"
	"sendSMS/comm"
)

var conf = flag.String("conf", "./etc/config.ini", "conf path")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始化配置文件
	err := config.InitModuleByParams(*conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("初始化配置文件")

	//初始化日志
	err = myLogger.InitLoggers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	myLogger.Infoln("初始化日志")

	//初始化db
	//err = gormdb.InitModule()
	//if err != nil {
	//	fmt.Println("初始化数据库失败", err)
	//	return
	//}
	//myLogger.Infoln("初始化数据库")
	//初始化缓存管理器
	//err = cache.InitModule()
	//if err != nil {
	//	fmt.Println("初始化缓存管理器失败", err)
	//	return
	//}
	//myLogger.Infoln("初始化缓存管理器")
	//初始化全局配置参数
	err = trans.InitArgv()
	if err != nil {
		fmt.Println("初始化全局配置参数失败", err)
		return
	}
	myLogger.Infoln("初始化全局配置参数")

	r := mux.NewRouter()
	r.HandleFunc("/sendSMS", handle.DoHandle)
	n := negroni.New()
	//n.Use(negroni.NewLogger())
	n.UseHandler(r)
	sockSvr := comm.HttpSvr{}
	err = sockSvr.InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sockSvr.RunSvr(n)

	os.Exit(0)
}
