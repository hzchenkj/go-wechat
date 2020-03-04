package wx

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/clbanning/mxj"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type weiXinQuery struct {
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string ` json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	EchoStr      string `json:"echostr"`
}

type WeiXinClient struct {
	Token          string
	Query          weiXinQuery
	Message        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Methods        map[string]func() bool
}

func NewClient(r *http.Request, w http.ResponseWriter, token string) (*WeiXinClient, error) {
	log.Println("New Client:" ,token)
	weiXinClient := new(WeiXinClient)
	weiXinClient.Token = token
	weiXinClient.Request = r
	weiXinClient.ResponseWriter = w

	weiXinClient.initWeiXinQuery()

	if weiXinClient.Query.Signature != weiXinClient.signature() {
		return nil, errors.New("Invalid Signature. ")
	}
	return weiXinClient, nil
}

func (client *WeiXinClient) initWeiXinQuery() {
	var q weiXinQuery
	q.Nonce = client.Request.URL.Query().Get("nonce")
	q.EchoStr = client.Request.URL.Query().Get("echostr")
	q.Signature = client.Request.URL.Query().Get("signature")
	q.Timestamp = client.Request.URL.Query().Get("timestamp")
	q.EncryptType = client.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = client.Request.URL.Query().Get("msg_signature")

	client.Query = q
}

func (client *WeiXinClient) signature() string {
	strSlice := sort.StringSlice{client.Token, client.Query.Timestamp, client.Query.Nonce}
	sort.Strings(strSlice)

	str := ""
	for _, s := range strSlice {
		str += s
	}

	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (client *WeiXinClient) initMessage() error {
	log.Println("initMessage in wx.go")
	body, err := ioutil.ReadAll(client.Request.Body)
	if err != nil {
		return err
	}

	m, err := mxj.NewMapXml(body)

	log.Println("in wx.go m:" ,m)
	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message. ")
	}

	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type. ")
	}

	client.Message = message

	log.Println(client.Message)

	return nil
}

func (client *WeiXinClient) text() {
	log.Println("text method in wx.go")
	inMsg, ok := client.Message["Content"].(string)

	if !ok {
		return
	}

	fmt.Println(inMsg)

	var reply TextMessage

	reply.InitBaseData(client, "text")
	reply.Content = value2CDATA(fmt.Sprintf("我收到的是：%s", inMsg))

	log.Println("content:" ,reply.Content)
	replyXml, err := xml.Marshal(reply)

	if err != nil {
		log.Println(err)
		client.ResponseWriter.WriteHeader(403)
		return
	}

	client.ResponseWriter.Header().Set("Content-Type", "text/xml")
	log.Println("replyXml:",replyXml)
	client.ResponseWriter.Write(replyXml)

}

func (client *WeiXinClient) Run() {
	log.Println("run method in wx.go")
	err := client.initMessage()
	log.Println("run method initMessage in wx.go")
	if err != nil {
		log.Println(err)
		client.ResponseWriter.WriteHeader(http.StatusForbidden)
		return
	}
	MsgType, ok := client.Message["MsgType"].(string)

	if !ok {
		client.ResponseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	switch MsgType {
	case "text":
		client.text()
		break
	default:
		break
	}
	return
}
