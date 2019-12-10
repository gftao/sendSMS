package handle

import (
	"net/http"
	"errors"
	"sendSMS/iServer"
	"sendSMS/trans"
)

func DoHandle(w http.ResponseWriter, r *http.Request) {
	iServer.IDoFunct(w, r, DoSvr)
}

func DoSvr(msg *trans.TransMessage) (iServer.IDoAppTrans, error) {

	switch msg.MsgBody.Tran_cd {
	case "10004003":
		return &iServer.T10004003{}, nil
	case "10003003":
		return &iServer.T10003003{}, nil
	default:
		return nil, errors.New("不识别的交易码: " + msg.MsgBody.Tran_cd)
	}
}
