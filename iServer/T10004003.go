package iServer

import (
	"sendSMS/trans"
	"encoding/xml"
	"mygolib/gerror"
	"mygolib/modules/run"
	"time"
	"mygolib/security"
	"net"
	"fmt"
	"io"
	"strconv"
	"errors"
	"github.com/axgle/mahonia"
)

type T10004003 struct {
	run.BaseWorker
	reqMsg  *trans.TransMessage
	rootReq *trans.Root
	rootRsp *trans.Root
	toRoot  []byte
	outRoot []byte
	en      mahonia.Encoder
	dec     mahonia.Decoder
}

func (t *T10004003) Init() gerror.IError {
	t.NodeName = "T10004003"
	t.en = mahonia.NewEncoder("GBK")
	if t.en == nil {
		return gerror.NewR(2009, nil, "非法的编码字符集")
	}
	t.dec = mahonia.NewDecoder("GBK")
	if t.dec == nil {
		return gerror.NewR(2009, nil, "非法的编码字符集")
	}
	return nil
}

func (t *T10004003) DoTrans(msg *trans.TransMessage) (gerr gerror.IError) {
	t.reqMsg = msg
	t.SetOrderId(msg.MsgBody.Order_id)

	gerr = t.BuildReq()
	if gerr != nil {
		return
	}
	gerr = t.toBack()
	if gerr != nil {
		return
	}
	gerr = t.AnalyRsp()
	if gerr != nil {
		return
	}

	return
}
func (t *T10004003) BuildReq() gerror.IError {
	if t.rootReq == nil {
		t.rootReq = new(trans.Root)
	}
	//t.rootReq.Envelope.Head.TxnCode = t.reqMsg.MsgBody.Tran_cd
	t.rootReq.Envelope.Head.TxnCode = "10002001"
	t.rootReq.Envelope.Head.Version = t.reqMsg.Version
	t.rootReq.Envelope.Head.Invoker = trans.GlobA.Invoker
	t.rootReq.Envelope.Head.UserName = t.reqMsg.MsgBody.User_name
	//t.rootReq.Envelope.Head.UserName = "test"
	t.rootReq.Envelope.Head.BranchId = t.reqMsg.MsgBody.Branch_id
	//t.rootReq.Envelope.Head.BranchId = "9999000001"
	t.rootReq.Envelope.Head.Charset = t.reqMsg.Encoding
	t.rootReq.Envelope.Head.Timestamp = time.Now().Format("20060102150405")
	t.rootReq.Envelope.Head.RequestId = t.reqMsg.MsgBody.Order_id
	t.rootReq.Envelope.Head.SignatureAlgorithm = "HmacMD5"

	t.rootReq.Envelope.Body.OrgCode = t.reqMsg.MsgBody.Ins_id_cd
	t.rootReq.Envelope.Body.MchtCd = t.reqMsg.MsgBody.Mcht_cd
	t.rootReq.Envelope.Body.MobilePhoneNo = t.reqMsg.MsgBody.Phone_no
	t.rootReq.Envelope.Body.MsgContent = t.reqMsg.MsgBody.Msg_conten
	t.rootReq.Envelope.Body.PlanTime = t.reqMsg.MsgBody.Send_time

	//e, err := xml.MarshalIndent(t.rootReq.Envelope, " ", "    ")
	e, err := xml.Marshal(t.rootReq.Envelope)
	if err != nil {
		return gerror.NewR(2009, err, "Marshal failed")
	}

	//t.rootReq.EnvelopeXML = t.en.ConvertString(string(e))
	t.rootReq.EnvelopeXML = string(e)
	t.Infof("签名串GBK编码:[%s]", t.rootReq.EnvelopeXML)
	t.SignReq()
	//r, err := xml.MarshalIndent(t.rootReq, " ", "    ")
	r, err := xml.Marshal(t.rootReq)
	if err != nil {
		return gerror.NewR(2009, err, "Marshal failed")
	}
	cx := xml.Header + string(r)
	cs := t.en.ConvertString(cx)
	t.toRoot = []byte(cs)
	t.Infof("请求串GBK编码:%s", cs)
	return nil
}

func (t *T10004003) AnalyRsp() gerror.IError {

	cs := t.dec.ConvertString(string(t.outRoot))
	t.Infof("UTF-8:%s", cs)

	err := xml.Unmarshal([]byte(cs), &t.rootRsp)
	if err != nil {
		return gerror.NewR(2009, err, "Unmarshal failed")
	}
	t.Infof("解析应答：%+v", *t.rootRsp)

	if t.rootRsp.Envelope.Head.ResCode != "00" {
		return gerror.New(2010, t.rootRsp.Envelope.Head.ResCode, err, t.rootRsp.Envelope.Head.ResMsg)
	}
	return nil
}

func (t *T10004003) SignReq() gerror.IError {
	h := security.HmacMd5(t.rootReq.EnvelopeXML, trans.GlobA.HmacKeyB)
	t.rootReq.Signature = h
	t.Info("签名成功：", t.rootReq.Signature)
	return nil
}

func (t *T10004003) toBack() gerror.IError {
	peerAddr, err := net.ResolveTCPAddr("tcp4", trans.GlobA.BackHost)
	if err != nil {
		return gerror.NewR(30041, err, "BindIP 非法值:[%s]", trans.GlobA.BackHost)
	}
	tcpConn, err := net.DialTCP("tcp4", nil, peerAddr)
	if err != nil {
		return gerror.NewR(30041, err, "DialTCP 失败:[%s]", peerAddr)
	}
	defer tcpConn.Close()
	tcpConn.SetDeadline(time.Now().Add(time.Second * time.Duration(10)))

	sndNum, err := tcpConn.Write([]byte(fmt.Sprintf("%04d", len(t.toRoot))))
	t.Info("tcpConn Write len: ", sndNum)
	if err != nil && sndNum != 4 {
		return gerror.NewR(30041, err, "发送报文长度失败:[%s]", t.reqMsg.MsgBody.Order_id)
	}
	totalNum := len(t.toRoot)
	for i := 0; i < totalNum; i += sndNum {
		sndNum, err = tcpConn.Write(t.toRoot[i:])
		if err != nil {
			return gerror.NewR(30041, err, "发送报文失败:[%s]", t.reqMsg.MsgBody.Order_id)
		}
	}
	lenBuf, err := ReadN(tcpConn, 4)
	if err != nil {

		return gerror.NewR(30041, err, "ReadN读取报文长度失败:[%s]", t.reqMsg.MsgBody.Order_id)
	}
	rcvNum, err := strconv.Atoi(string(lenBuf))
	if err != nil || rcvNum > 9999 {
		return gerror.NewR(30041, err, "ReadN读取报文长度失败:[%s]", t.reqMsg.MsgBody.Order_id)
	}
	rp, err := ReadN(tcpConn, rcvNum)
	if err != nil {
		return gerror.NewR(30041, err, "ReadN读取响应报文非法:[%s]", t.reqMsg.MsgBody.Order_id)
	}
	//rsp := string(rp)[len(xml.Header)-1:]
	t.Infof("收到应答报文[%s]", rp)
	t.outRoot = rp
	return nil
}

func ReadN(rd io.Reader, len int) ([]byte, error) {
	buf := make([]byte, len)
	rcvNum := 0
	totalNum := 0
	var err error
	for i := 0; i < len; i += rcvNum {
		rcvNum, err = rd.Read(buf[i:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		totalNum += rcvNum
		if err == io.EOF {
			break
		}
	}
	if len != totalNum {
		return nil, errors.New("读取长度失败")
	}
	return buf, nil
}
