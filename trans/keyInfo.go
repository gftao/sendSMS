package trans

import (
	"crypto/rsa"
	"mygolib/security"
	"mygolib/gerror"
	"mygolib/defs"
	"mygolib/modules/cache"
	"mygolib/modules/gormdb"
	"time"
	"upEletrcSign/dbModels"
	"mygolib/modules/myLogger"
)

type KeyInfo struct {
	Key_down_flg   string `json:"key_down_flg,omitempty"`
	Session_key    string `json:"session_key,omitempty"`
	Term_pub_key   string `json:"term_pub_key,omitempty"`
	Server_pub_key string `json:"server_pub_key,omitempty"`
	Term_tmk       string `json:"term_tmk,omitempty"`
	Term_tmk_chk   string `json:"term_tmk_chk,omitempty"`
	Term_pik       string `json:"term_pik,omitempty"`
	Term_pik_chk   string `json:"term_pik_chk,omitempty"`
	Term_trk       string `json:"term_trk,omitempty"`
	Term_trk_chk   string `json:"term_trk_chk,omitempty"`
	Pboc_ic_key    string `json:"pboc_ic_key,omitempty"`
	Pboc_ic_param  string `json:"pboc_ic_param,omitempty"`
}

type KeyHandleInfo struct {
	Term_key     string
	TermPubKey   *rsa.PublicKey
	TermPriKey   *rsa.PrivateKey
	ServerPubKey *rsa.PublicKey
	ServerPriKey *rsa.PrivateKey
}

func InitKeyHandle(k *KeyHandleInfo, prefix, mcht_cd, term_id string) gerror.IError {
	//从cache中查找句柄
	term_key := prefix + mcht_cd + term_id
	err := cache.Get(term_key, k)
	if err == cache.ErrCacheMiss {
		//取出终端密钥信息
		termInfo := dbModels.DbTermInfo{}
		dbc := gormdb.GetInstance()
		err := dbc.Where("mcht_cd = ? and term_id = ? ",
			mcht_cd, term_id).Find(&termInfo).Error
		if err != nil {
			return gerror.New(10040, defs.TRN_SYS_ERROR, err, "取终端密钥信息失败")
		}

		myLogger.Info(term_key + "密钥句柄未找到,重新加载")
		prikey, err := security.GetRsaPrivateKeyByString(termInfo.ServerPriKey)
		if err != nil {
			return gerror.New(10050, defs.TRN_SYS_ERROR, err, "取终端密钥信息失败")
		}
		k.ServerPriKey = prikey
		pubkey, err := security.GetRsaPublicKeyByString(termInfo.TermPubKey)
		if err != nil {
			return gerror.New(10060, defs.TRN_SYS_ERROR, err, "取终端密钥信息失败")
		}
		k.TermPubKey = pubkey
		//保存key到缓存
		err = cache.Add(term_key, *k, time.Duration(TermKeyOutTime)*time.Second)
		if err != nil {
			return gerror.New(10070, defs.TRN_SYS_ERROR, err, "保存密钥句柄到缓存失败")
		}
	} else if err != nil {
		return gerror.New(10080, defs.TRN_SYS_ERROR, err, "取密钥缓存出错")
	}
	myLogger.Info("从缓存加载key成功:" + term_key)
	return nil
}
