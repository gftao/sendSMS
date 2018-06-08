package iServer

import (
	"net/http"
	"sendSMS/trans"
	"io/ioutil"
	"mygolib/modules/myLogger"
	"mygolib/defs"
	"mygolib/gerror"
)

type GetDoTransFunc func(msg *trans.TransMessage) (IDoAppTrans, error)

type IDoAppTrans interface {
	Init() gerror.IError
	DoTrans(*trans.TransMessage) (gerror.IError)
}

func IDoFunct(w http.ResponseWriter, r *http.Request, getTransFunc GetDoTransFunc) {
	ra, err := ioutil.ReadAll(r.Body)
	if err != nil {
		myLogger.Error("读取请求报文失败", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	myLogger.Infof("get Request msg:[%s]", string(ra))
	tr, gerr := trans.UnPackReq(ra)
	if gerr != nil {
		myLogger.Errorf("解析请求失败:[%s]", gerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	gerr = trans.VerifyTransMessage(tr)
	if gerr != nil {
		myLogger.Errorf("报文验证失败:[%s]", gerr)
		RejectMsg(w, tr, defs.TRN_FORMAT_ERR, err.Error())
		return
	}

	TransFunc, err := getTransFunc(tr)
	if err != nil {
		RejectMsg(w, tr, defs.TRN_FORMAT_ERR, err.Error())
		return
	}
	gerr = TransFunc.Init()
	if gerr != nil {
		RejectMsg(w, tr, defs.TRN_FORMAT_ERR, err.Error())
		return
	}
	gerr = TransFunc.DoTrans(tr)
	if gerr != nil {
		myLogger.Error("交易处理失败", gerr)
		RejectMsg(w, tr, gerr.GetErrorCode(), gerr.GetErrorString())
		return
	}
	myLogger.Info("应答处理完成")

	tr.MsgBody.Resp_cd = "00"
	tr.MsgBody.Resp_msg = "SUCCESS"
	gerr = trans.SignTransMessage(tr)
	if gerr != nil {
		myLogger.Errorf("报文签名失败:[%s]", gerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rep := tr.ToString()
	myLogger.Infof("应答报文：[%s]", rep)

	w.Write([]byte(rep))
	return
}

func RejectMsg(w http.ResponseWriter, msg *trans.TransMessage, resp_cd, resp_msg string) {

	myLogger.Debugf("resp_cd:%s, resp_msg:%s", resp_cd, resp_msg)

	msg.MsgBody.Resp_cd = resp_cd
	msg.MsgBody.Resp_msg = resp_msg
	err := trans.SignTransMessage(msg)
	if err != nil {
		myLogger.Errorf("报文签名失败:[%s]", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rep := msg.ToString()
	myLogger.Debug("本地拒绝成功: ", string(rep))
	w.Write([]byte(rep))
}
