package iServer

import (
	"sendSMS/trans"
	"mygolib/gerror"
	"mygolib/modules/run"
	"github.com/axgle/mahonia"
	"sendSMS/models"
	"mygolib/modules/gormdb"
)

type T10003003 struct {
	run.BaseWorker
	reqMsg  *trans.TransMessage
	rootReq *trans.Root
	rootRsp *trans.Root
	toRoot  []byte
	outRoot []byte
	en      mahonia.Encoder
	dec     mahonia.Decoder
}

func (t *T10003003) Init() gerror.IError {

	return nil
}

func (t *T10003003) DoTrans(msg *trans.TransMessage) (gerr gerror.IError) {
	t.reqMsg = msg
	t.SetOrderId(msg.MsgBody.Mcht_cd)
	if msg.MsgBody.Mcht_cd == "" {
		return gerror.New(80112, "32", nil, "无效商户号%s", t.reqMsg.MsgBody.Mcht_cd)
	}
	mcht := models.T_m_mcht_inf{}
	dbc := gormdb.GetInstance()
	err := dbc.Where("mcht_cd = ?", t.reqMsg.MsgBody.Mcht_cd).Find(&mcht).Error
	if err != nil {
		return gerror.New(80111, "33", err, "无效商户号%s", t.reqMsg.MsgBody.Mcht_cd)
	}
	if mcht.Mobile == "" {
		return gerror.New(80113, "34", err, "商户%s未录入手机号", t.reqMsg.MsgBody.Mcht_cd)
	}
	t.reqMsg.MsgBody.Phone_no = mcht.Mobile
	return nil
}
