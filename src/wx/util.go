package wx

import (
	"encoding/xml"
	"strconv"
	"time"
)



type Base struct {
	FromUsername CDATAText
	ToUsername CDATAText
	MsgType CDATAText
	CreateTime CDATAText
}

func (b *Base) InitBaseData(w *WeiXinClient, msgType string){
	b.FromUsername = value2CDATA(w.Message["ToUserName"].(string))
	b.ToUsername = value2CDATA(w.Message["FromUserName"].(string))
	b.CreateTime = value2CDATA(strconv.FormatInt(time.Now().Unix(),10))
	b.MsgType = value2CDATA(msgType)
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

type TextMessage struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Content CDATAText
}

func value2CDATA(v string) CDATAText {
	return  CDATAText{"<![CDATA[" + v + "]]>"}
}