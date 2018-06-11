package trans

import (
	"crypto/rsa"
	"mygolib/modules/config"
	"errors"
	"mygolib/modules/myLogger"
)

const TermKeyOutTime = 100

type globArgv struct {
	HmacKeyS string
	HmacKeyB string
	BackHost string
	Invoker  string
	PubKey   *rsa.PublicKey
	PriKey   *rsa.PrivateKey
}

var GlobA *globArgv

func InitArgv() error {
	if !config.HasConfigInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}

	GlobA = new(globArgv)
	config.SetSection("glob")

	GlobA.HmacKeyS = config.StringDefault("HmacKeyS", "")
	GlobA.HmacKeyB = config.StringDefault("HmacKeyB", "")
	ip := config.StringDefault("RemoteIP", "")
	port := config.StringDefault("RemptePort", "")
	GlobA.BackHost = ip + ":" + port
	GlobA.Invoker = config.StringDefault("Invoker", "")

	myLogger.Infof("GlobA:%+v", *GlobA)

	return nil
}
