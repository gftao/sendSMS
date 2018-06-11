package trans

import (
	"encoding/json"
	"bytes"
	"compress/zlib"
	"io"
	"encoding/base64"
	"compress/gzip"
	"mygolib/modules/myLogger"
	"mygolib/security"
	"mygolib/defs"
	"mygolib/gerror"
)

func UnPackReq(req []byte) (*TransMessage, gerror.IError) {
	tranMsg := TransMessage{}

	err := json.Unmarshal(req, &tranMsg)
	if err != nil {
		return nil, gerror.NewR(1001, err, "解析请求失败")
	}
	myLogger.Infof("get Msg_body:%s", tranMsg.Msg_body)
	err = json.Unmarshal([]byte(tranMsg.Msg_body), &tranMsg.MsgBody)
	if err != nil {
		return nil, gerror.NewR(1002, err, "解析msgbody失败")
	}
	if tranMsg.MsgBody == nil {
		return nil, gerror.NewR(1003, nil, "msg_body is nil")
	}

	return &tranMsg, nil
}
func Md5Base64(origData []byte) string {
	return security.GenMd5(origData)
}
func SignTransMessage(t *TransMessage) gerror.IError {

	if t.Sign_method == "01" {
		myLogger.Debug("开始RSA签名")

	} else if t.Sign_method == "03" {
		myLogger.Debug("开始HmacMD5签名")
		h := security.HmacMd5(t.Msg_body, GlobA.HmacKeyS)
		t.Signature = h
	}
	myLogger.Debug("报文签名成功")
	return nil
}

func VerifyTransMessage(t *TransMessage) gerror.IError {
	var ok bool = false
	if t.Sign_method == "01" {
		myLogger.Debug("开始RSA验证")

	} else if t.Sign_method == "02" {
		myLogger.Debug("开始MD5验证")
		ok = security.VerifyMd5([]byte(t.Msg_body), t.Signature)
	} else if t.Sign_method == "03" {
		myLogger.Debug("开始HmacMD5验证")
		ok = security.VerifyHmacMd5([]byte(t.Msg_body), []byte(t.Signature), []byte(GlobA.HmacKeyS))
	} else if t.Sign_method == "AA" {
		myLogger.Debug("测试不验证签名")
		ok = true
	}
	if !ok {
		return gerror.New(10070, defs.TRN_VERIFY_MAC_FAIL, nil, "报文验证不通过")
	}
	myLogger.Debug("验证报文成功")
	return nil

}

func DecodeBase64(cipherdata []byte) (string, error) {
	orig, err := base64.StdEncoding.DecodeString(string(cipherdata))

	return string(orig), err
}

func UnDoZlibCompressBase64(src string) (string, error) {
	us, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return UnDoZlibCompress(us)
}

func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func UnDoZlibCompress(src []byte) (string, error) {
	var out bytes.Buffer

	i := bytes.NewReader(src)
	r, err := zlib.NewReader(i)
	if err != nil {
		return "", err
	}
	io.Copy(&out, r)
	r.Close()

	return out.String(), err
}

func UnDoGzipCompressBase64(src string) (string, error) {
	us, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return UnDoGzipCompress(us)
}

func UnDoGzipCompress(src []byte) (string, error) {
	var out bytes.Buffer

	i := bytes.NewReader(src)
	r, err := gzip.NewReader(i)
	if err != nil {
		return "", err
	}
	io.Copy(&out, r)
	r.Close()

	return out.String(), err
}
