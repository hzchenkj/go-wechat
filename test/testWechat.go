package main

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"test"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	argNum := len(os.Args)
	if argNum < 6 {
		fmt.Println("No Enough Args")
		return
	}

	url := os.Args[1]
	token := os.Args[2]
	from := os.Args[3]
	to := os.Args[4]
	text := os.Args[5]

	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	nonce := test.RandStringRunes(8)

	sign := test.Signature(timestampStr, nonce, token)

	url = fmt.Sprintf("%s?signature=%s&timestamp=%s&nonce=%s", url, sign, timestampStr, nonce)

	var message test.Text
	message.FromUserName = test.StrToCDATA(from)
	message.ToUserName = test.StrToCDATA(to)
	message.MsgType = test.StrToCDATA("text")
	message.MsgId = rand.Int63()
	message.Content = test.StrToCDATA(text)
	message.CreateTime = timestamp

	xmlMsg, err := xml.Marshal(message)

	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := test.Send(url, string(xmlMsg))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("URL:", url)
	fmt.Println("--------------------")
	fmt.Println("Send Message:")

	x, err := test.FormatXML(xmlMsg)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(x))
	}

	fmt.Println("--------------------")
	fmt.Println("Response:")

	if resp == nil {
		fmt.Println("--------------------")
		return
	}

	x, err = test.FormatXML(resp)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(x))
	}

	fmt.Println("--------------------")
}