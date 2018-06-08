package trans

import (
	"encoding/json"
)

type MessageHead struct {
	Encoding    string `json:"encoding"`
	Sign_method string `json:"sign_method"`
	Signature   string `json:"signature"`
	Version     string `json:"version"`
}

type TransMessage struct {
	MessageHead
	Msg_body string       `json:"msg_body"`
	MsgBody  *TransParams `json:"-"`
}

type TransParams struct {
	commonParams
}

type commonParams struct {
	Tran_cd    string `json:"tran_cd,omitempty"`
	Mcht_cd    string `json:"mcht_cd,omitempty"`
	Resp_cd    string `json:"resp_cd,omitempty"`
	Resp_msg   string `json:"resp_msg,omitempty"`
	Ins_id_cd  string `json:"ins_id_cd,omitempty"`
	Send_time  string `json:"send_time,omitempty"`
	Order_id   string `json:"order_id,omitempty"`
	Phone_no   string `json:"phone_no,omitempty"`
	Msg_conten string `json:"msg_conten,omitempty"`
	User_name  string `json:"user_name,omitempty"`
	Branch_id  string `json:"branch_id,omitempty"`
}

func (t *TransMessage) SetMsgBody() {
	btMsgBody, err := json.Marshal(t.MsgBody)
	if err != nil {
		t.Msg_body = "{}"
		return
	}
	t.Msg_body = string(btMsgBody)
}

func (t *TransMessage) ToString() string {
	t.SetMsgBody()
	res, err := json.Marshal(t)
	if err != nil {
		return "{}"
	}
	return string(res)
}

type Root struct {
	XMLName     struct{} `xml:"root"`
	Signature   string   `xml:"signature,omitempty"`
	EnvelopeXML string   `xml:"-"`
	Envelope    Envelope
}
type Envelope struct {
	XMLName struct{} `xml:"envelope"`
	Head    Head     `xml:"head,omitempty"`
	Body    Body     `xml:"body,omitempty"`
}
type Head struct {
	Key                string `xml:"key,omitempty"`
	Version            string `xml:"version,omitempty"`
	TxnCode            string `xml:"txnCode,omitempty"`
	Invoker            string `xml:"invoker,omitempty"`
	UserName           string `xml:"userName,omitempty"`
	BranchId           string `xml:"branchId,omitempty"`
	Charset            string `xml:"charset,omitempty"`
	Timestamp          string `xml:"timestamp,omitempty"`
	RequestId          string `xml:"requestId,omitempty"`
	SignatureAlgorithm string `xml:"signatureAlgorithm,omitempty"`
	ResCode            string `xml:"resCode,omitempty"`
	ResMsg             string `xml:"resMsg,omitempty"`
}
type Body struct {
	OrgCode       string `xml:"orgCode,omitempty"`
	MchtCd        string `xml:"mchtCd,omitempty"`
	MobilePhoneNo string `xml:"mobilePhoneNo,omitempty"`
	MsgContent    string `xml:"msgContent,omitempty"`
	PlanTime      string `xml:"planTime,omitempty"`
}
